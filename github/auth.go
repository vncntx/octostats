package github

import (
	"encoding/base64"
	"fmt"
)

// Credentials are Basic Authentication credentials for GitHub
// See https://docs.github.com/en/rest/overview/other-authentication-methods
type Credentials struct {
	User  string
	Token string
}

// EncodeToBasic encodes the credentials as an RFC-2617 Basic Authorization header
func (c Credentials) EncodeToBasic() string {
	credentials := fmt.Sprintf("%s:%s", c.User, c.Token)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}
