package main

import (
	"os"

	"github.com/vincentfiestada/octostats/github"
)

func main() {
	auth := github.Credentials{
		User:  os.Getenv(GitHubUser),
		Token: os.Getenv(GitHubToken),
	}
	if len(auth.User) < 1 {
		log.Error("GitHub login user is empty. Set %s env variable", GitHubUser)
	}
	if len(auth.Token) < 1 {
		log.Error("GitHub login token is empty. Set %s env variable", GitHubToken)
	}

	client := github.NewOctoClient(auth)

	user, err := client.GetAuthenticatedUser()
	if err != nil {
		log.Fatal("failed to get authenticated user: %s", err)
		return
	}
	log.Info("authenticated as %#v", user)
}
