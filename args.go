package main

import (
	"os"
	"strings"
	"time"

	"github.com/vincentfiestada/octostats/commons"
	"github.com/vincentfiestada/octostats/github"
	"github.com/vincentfiestada/octostats/util"
)

// getAuth returns credentials from environment variables
func getAuth() github.Credentials {
	auth := github.Credentials{
		User:  os.Getenv(GitHubUser),
		Token: os.Getenv(GitHubToken),
	}

	if len(auth.User) < 1 {
		log.Fatal("GitHub login user is empty. Set %s env variable", GitHubUser)
	}
	if len(auth.Token) < 1 {
		log.Fatal("GitHub login token is empty. Set %s env variable", GitHubToken)
	}

	return auth
}

// getRepo returns the repository name from command-line arguments
func getRepo() string {
	if len(os.Args) < 2 {
		log.Fatal("repository name is required.")
	}
	repo := os.Args[1]
	if !util.Matches(repositoryPattern, repo) {
		log.Fatal("repository name '%s' is invalid; must be of the form `owner/name`", repo)
	}

	return repo
}

// getStartTime returns the time to start collecting stats from command-line arguments
func getStartTime() time.Time {
	t := time.Time{}
	if len(os.Args) >= 3 {
		timeArg := os.Args[2]
		if strings.HasPrefix(timeArg, "-") {
			// use argument as time offset
			t = getTimeFromDuration(timeArg)
		} else {
			// use argument as date
			t = getTimeFromDate(timeArg)
		}
	}

	return t
}

func getTimeFromDuration(rawDuration string) time.Time {
	duration, err := time.ParseDuration(rawDuration)
	if err != nil {
		log.Fatal("invalid duration '%s'", rawDuration)
	}

	return time.Now().Add(duration)
}

func getTimeFromDate(rawDate string) time.Time {
	t, err := time.Parse(commons.DateLayout, rawDate)
	if err != nil {
		log.Fatal("invalid date '%s'", rawDate)
	}

	return t
}
