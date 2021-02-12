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
		log.Fatal("GitHub login user is empty. Set %s env variable", GitHubUser)
		return
	}
	if len(auth.Token) < 1 {
		log.Fatal("GitHub login token is empty. Set %s env variable", GitHubToken)
		return
	}

	if len(os.Args) < 3 {
		log.Fatal("not enough arguments")
		return
	}
	owner := os.Args[1]
	repo := os.Args[2]

	client := github.NewOctoClient(auth)

	user, err := client.GetAuthenticatedUser()
	if err != nil {
		log.Fatal("failed to get authenticated user: %s", err)
		return
	}
	log.Info("authenticated as %#v", user)
	log.Info("inspecting %s/%s", owner, repo)
}
