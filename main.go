package main

import (
	"os"

	"github.com/vincentfiestada/octostats/github"
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
	count := 0 // number of merged pull requests
	// totalInNanoseconds := int64(0) // total time to merge in nanoseconds

	query := github.Query("").WithRepo(repo).WithAuthor(auth.User).IsMerged()

	for {
		log.Debug("getting page %d of pull requests", page)
		results, err := client.SearchPulls(repo, page, query)
		if err != nil {
			log.Fatal("failed to get pull requests: %s", err)
			return
		}
		if results.IsIncomplete {
			log.Warn("search results are incomplete due to timeout")
		}
		if len(results.Items) < 1 {
			break
		}
		page++

		for _, pull := range results.Items {
			log.Debug("%#v", pull)

			// timeToMerge := pull.MergedAt.Sub(pull.CreatedAt)
			// log.Info("pull request #%d took %.6f hours to merge (created by %s)", pull.Number, timeToMerge.Hours(), pull.User)

			// totalInNanoseconds += timeToMerge.Nanoseconds()
			count++
		}
	}

	log.Info("found %d merged pull requests for %s", count, repo)

	// if count > 0 {
	// 	avg := time.Duration(totalInNanoseconds / int64(count))
	// 	log.Info("average time to merge: %.6f hours", avg.Hours())
	// }
}
