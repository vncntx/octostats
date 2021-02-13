package github

// GitHub endpoints
const (
	baseEndpoint = "https://api.github.com"

	userEndpoint   = "/user"
	reposEndpoint  = "/repos"
	pullsEndpoint  = reposEndpoint + "/%s/pulls"
	searchEndpoint = "/search/issues"
)

// GitHub search qualifiers
const (
	pr     = "pr"
	draft  = "draft"
	merged = "merged"
)

// GitHub sort methods
const (
	updated = "updated"
)

// GitHub sort ordering
const (
	descending = "desc"
)
