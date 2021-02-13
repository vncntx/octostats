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
	reviewsCount := 0                          // total number of pr reviews
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

			// Record time to merge
			if pull.CreatedAt.IsZero() || pull.MergedAt.IsZero() {
				log.Warn("pull request #%d has invalid timestamps", pull.Number)
				continue
			}

			timeToMerge := pull.MergedAt.Sub(pull.CreatedAt)
			timeToMergeInNanoSeconds += timeToMerge.Nanoseconds()

			// Record number of reviews
			reviews, err := countReviews(client, repo, item.Number)
			if err != nil {
				log.Warn("couldn't count reviews for pull request #%d: %s", item.Number, err)
				continue
			}
			reviewsCount += reviews

			log.Info("pull request #%d took %.6f hours to merge (%d reviews)", pull.Number, timeToMerge.Hours(), reviews)

			// Count per label
			for _, label := range pull.Labels {
				countByLabel[label]++
			}

			// Count merged pull requests
			count++
		}
	}

	log.Info("found %d merged pull requests for %s", count, repo)

	if count < 1 {
		return
	}

	avgTimeToMerge := time.Duration(timeToMergeInNanoSeconds / int64(count))
	log.Info("average time to merge: %.6f hours", avgTimeToMerge.Hours())

	avgReviewsCount := reviewsCount / count
	log.Info("average reviewers count: %d", avgReviewsCount)

	for label, count := range countByLabel {
		log.Info("with label '%s': %d", label, count)
	}
}

func countReviews(client github.Client, repo string, pr int) (int, error) {
	page := 1
	count := 0

	for {
		reviews, err := client.ListReviews(repo, pr, page)
		if err != nil {
			log.Fatal("failed to get pull requests: %s", err)
			return count, err
		}
		if len(reviews) < 1 {
			break
		}
		page++

		count += len(reviews)
	}

	return count, nil
}
