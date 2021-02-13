package github

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Client wraps the GitHub API
type Client interface {
	GetAuthenticatedUser() (User, error)
	ListPulls(repo string, page int, state string) ([]Pull, error)
	SearchPulls(repo string, page int, q *QueryBuilder) (SearchResponse, error)
	GetPull(repo string, number int) (Pull, error)
}

// User is a GitHub user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Email string `json:"email"`
}

func (u User) String() string {
	return u.Login
}

// Pull is a GitHub Pull Request
type Pull struct {
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	User        User      `json:"user"`
	IsMerged    bool      `json:"merged"`
	CreatedAt   time.Time `json:"created_at"`
	MergedAt    time.Time `json:"merged_at"`
	ClosedAt    time.Time `json:"closed_at"`
	MergeCommit string    `json:"merge_commit_sha"`
	Reviewers   []User    `json:"requested_reviewers"`
}

// SearchResponse are the results of an issue or pull request search
type SearchResponse struct {
	Total        int            `json:"total_count"`
	IsIncomplete bool           `json:"incomplete_results"`
	Items        []SearchResult `json:"items"`
}

// SearchResult is a single issue or pull request search result
type SearchResult struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	User      User      `json:"user"`
	State     string    `json:"state"`
	Repo      Repo      `json:"repository_url"`
	CreatedAt time.Time `json:"created_at"`
}

// Repo is the unique name of a GitHub repository
type Repo string

// UnmarshalJSON parses a repo name from a url
func (r *Repo) UnmarshalJSON(b []byte) error {
	var url string
	if err := json.Unmarshal(b, &url); err != nil {
		return err
	}

	prefix := baseEndpoint + reposEndpoint + "/"
	if !strings.HasPrefix(url, prefix) {
		return fmt.Errorf("unable to parse invalid repo URL '%s'", url)
	}

	name := strings.TrimPrefix(url, prefix)
	*r = Repo(name)
	return nil
}
