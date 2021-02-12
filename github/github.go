package github

// Client wraps the GitHub API
type Client interface {
	GetAuthenticatedUser() (User, error)
}

// User is a GitHub user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Email string `json:"email"`
}
