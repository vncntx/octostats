package github

// GitHub endpoints
const (
	BaseEndpoint = "https://api.github.com"

	UserEndpoint  = "/user"
	ReposEndpoint = "/repos/%s"
	PullsEndpoint = ReposEndpoint + "/pulls"
)
