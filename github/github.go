package github

import "time"

// Client wraps the GitHub API
type Client interface {
	GetAuthenticatedUser() (User, error)
	ListPulls(repo string, page int, state string) ([]Pull, error)
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
	CreatedAt   time.Time `json:"created_at"`
	MergedAt    time.Time `json:"merged_at"`
	MergeCommit string    `json:"merge_commit_sha"`
}
