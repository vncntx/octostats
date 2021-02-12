package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		BaseEndpoint: BaseEndpoint,

		auth:       cred.EncodeToBasic(),
		httpClient: &http.Client{},
	}
}

// GetAuthenticatedUser returns the current user
func (o *OctoClient) GetAuthenticatedUser() (User, error) {
	res, err := o.getRequest(UserEndpoint, nil)
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

func (o *OctoClient) getRequest(endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", o.buildURL(endpoint), body)
	if err != nil {
		return nil, err
	}
	o.setHeaders(req)
	return o.sendRequest(req)
}

func (o *OctoClient) buildURL(endpoint string) string {
	return o.BaseEndpoint + endpoint
}

func (o *OctoClient) setHeaders(req *http.Request) {
	req.Header.Add("Authorization", o.auth)
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
