package main

import (
	"os"
	"time"

	"github.com/vincentfiestada/octostats/github"
	"github.com/vincentfiestada/octostats/github/filters"
	"github.com/vincentfiestada/octostats/util"
)

const (
	repositoryPattern = "[A-Za-z0-9-_]+/[A-Za-z0-9-_]+"
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

	if len(os.Args) < 2 {
		log.Fatal("repository name is required.")
		return
	}
	repo := os.Args[1]
	if !util.Matches(repositoryPattern, repo) {
		log.Fatal("repository name '%s' is invalid; must be of the form `owner/name`", repo)
		return
	}
	log.Info("inspecting %s", repo)

	client := github.NewOctoClient(auth)

	if user, err := client.GetAuthenticatedUser(); err != nil {
		log.Error("failed to get authenticated user: %s", err)
	} else {
		log.Info("authenticated as %s", user)
	}

	page := 1
	count := 0                     // number of merged pull requests
	totalInNanoseconds := int64(0) // total time to merge in nanoseconds
	for {
		log.Debug("getting page %d of pull requests", page)
		pulls, err := client.ListPulls(repo, page, filters.All)
		if err != nil {
			log.Fatal("failed to get pull requests user: %s", err)
			return
		}
		if len(pulls) < 1 {
			break
		}
		page++

		for _, pull := range pulls {
			if !pull.MergedAt.IsZero() && pull.User.Login == auth.User {
				timeToMerge := pull.MergedAt.Sub(pull.CreatedAt)
				log.Info("pull request #%d took %.6f hours to merge (created by %s)", pull.Number, timeToMerge.Hours(), pull.User)

				totalInNanoseconds += timeToMerge.Nanoseconds()
				count++
			}
		}
	}

	avg := time.Duration(totalInNanoseconds / int64(count))
	log.Info("found %d merged pull requests for %s", count, repo)
	log.Info("average time to merge: %.6f hours", avg.Hours())
}
