package main

import (
	"os"
	"time"

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
	count := 0                                 // number of merged pull requests
	countByLabel := make(map[github.Label]int) // number of prs by label
	reviewersCount := 0                        // total number of reviewers
	timeToMergeInNanoSeconds := int64(0)       // total time to merge in nanoseconds

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

		for _, item := range results.Items {

			pull, err := client.GetPull(repo, item.Number)
			if err != nil {
				log.Warn("couldn't get details of pull request #%d: %s", item.Number, err)
				continue
			}
			if pull.CreatedAt.IsZero() || pull.MergedAt.IsZero() {
				log.Warn("pull request #%d has invalid timestamps", pull.Number)
				continue
			}

			timeToMerge := pull.MergedAt.Sub(pull.CreatedAt)
			log.Info("pull request #%d took %.6f hours to merge (%d reviewers)", pull.Number, timeToMerge.Hours(), len(pull.Reviewers))

			reviewersCount += len(pull.Reviewers)
			timeToMergeInNanoSeconds += timeToMerge.Nanoseconds()
			count++

			for _, label := range pull.Labels {
				countByLabel[label]++
			}
		}
	}

	log.Info("found %d merged pull requests for %s", count, repo)

	if count < 1 {
		return
	}

	avgTimeToMerge := time.Duration(timeToMergeInNanoSeconds / int64(count))
	log.Info("average time to merge: %.6f hours", avgTimeToMerge.Hours())

	avgReviewersCount := reviewersCount / count
	log.Info("average reviewers count: %d", avgReviewersCount)

	for label, count := range countByLabel {
		log.Info("with label '%s': %d", label, count)
	}
}
