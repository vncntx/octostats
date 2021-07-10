package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// OctoClient is a GitHub REST API Client
type OctoClient struct {
	BaseEndpoint string

	auth       string
	httpClient *http.Client
}

// NewOctoClient returns a new OctoClient
func NewOctoClient(cred Credentials) Client {
	return &OctoClient{
		BaseEndpoint: baseEndpoint,

		auth:       cred.EncodeToBasic(),
		httpClient: &http.Client{},
	}
}

// GetAuthenticatedUser returns the current user
func (o *OctoClient) GetAuthenticatedUser() (User, error) {
	res, err := o.sendGet(userEndpoint, nil, nil)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = o.parseBody(res, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// ListPulls returns the list of pull requests
func (o *OctoClient) ListPulls(repo string, page int, state string) ([]Pull, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("state", state)

	endpoint := fmt.Sprintf(pullsEndpoint, repo)

	res, err := o.sendGet(endpoint, params, nil)
	if err != nil {
		return nil, err
	}

	pulls := []Pull{}
	err = o.parseBody(res, &pulls)
	if err != nil {
		return nil, err
	}

	return pulls, nil
}

// SearchPulls searches for pull requests
func (o *OctoClient) SearchPulls(repo string, page int, q *QueryBuilder) (SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("sort", updated)
	params.Set("order", descending)
	params.Set("q", q.withType(pr).Build())

	res, err := o.sendGet(searchEndpoint, params, nil)
	if err != nil {
		return SearchResponse{}, err
	}

	results := SearchResponse{}
	err = o.parseBody(res, &results)
	if err != nil {
		return SearchResponse{}, err
	}

	return results, nil
}

// GetPull retrieves details about a pull request
func (o *OctoClient) GetPull(repo string, pr int) (Pull, error) {
	endpoint := fmt.Sprintf(pullDetailsEndpoint, repo, pr)

	res, err := o.sendGet(endpoint, nil, nil)
	if err != nil {
		return Pull{}, err
	}

	pull := Pull{}
	err = o.parseBody(res, &pull)
	if err != nil {
		return Pull{}, err
	}

	return pull, nil
}

// ListReviews returns the list of reviews for a pull request
func (o *OctoClient) ListReviews(repo string, pr int, page int) ([]Review, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))

	endpoint := fmt.Sprintf(reviewsEndpoint, repo, pr)

	res, err := o.sendGet(endpoint, params, nil)
	if err != nil {
		return nil, err
	}

	reviews := []Review{}
	err = o.parseBody(res, &reviews)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// sendGet sends an HTTP GET request
func (o *OctoClient) sendGet(endpoint string, params url.Values, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", o.buildURL(endpoint, params), body)
	if err != nil {
		return nil, err
	}
	o.setHeaders(req)

	return o.sendRequest(req)
}

func (o *OctoClient) buildURL(endpoint string, params url.Values) string {
	url := o.BaseEndpoint + endpoint
	if params != nil {
		url += "?" + params.Encode()
	}

	return url
}

func (o *OctoClient) setHeaders(req *http.Request) {
	req.Header.Add("Authorization", o.auth)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
}

func (o *OctoClient) sendRequest(req *http.Request) (*http.Response, error) {
	res, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("HTTP %d received from GitHub", res.StatusCode)
	}

	return res, nil
}

func (o *OctoClient) parseBody(res *http.Response, v interface{}) error {
	if res.Body == nil {
		return fmt.Errorf("no response from GitHub: %s", res.Request.URL)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}
