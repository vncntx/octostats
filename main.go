package main

import (
	"time"

	"octostats/commons"
	"octostats/github"

	"github.com/hako/durafmt"
)

const (
	repositoryPattern = "[A-Za-z0-9-_]+/[A-Za-z0-9-_]+"
)

func main() {
	auth := getAuth()
	repo := getRepo()
	mergedAfter := getStartTime()

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

	query := buildQuery(repo, auth, mergedAfter)

	for {
		results := runQuery(client, repo, query, page)
		if len(results.Items) < 1 {
			break
		}
		page++

		for _, item := range results.Items {
			stats, err := getStats(client, repo, item)
			if err != nil {
				log.Warn(err.Error())

				continue
			}
			stats.log()

			// aggregate stats
			timeToMergeInNanoSeconds += stats.timeToMerge.Nanoseconds()
			reviewsCount += stats.reviews
			for _, label := range stats.labels {
				countByLabel[label]++
			}
			count++
		}
	}

	avgTimeToMerge := time.Duration(timeToMergeInNanoSeconds / int64(count))
	avgReviews := reviewsCount / count

	summarize(repo, count, avgTimeToMerge, avgReviews, countByLabel)
}

// buildQuery builds a query for pull requests by repo, author, and merge time
func buildQuery(repo string, auth github.Credentials, mergedAfter time.Time) *github.QueryBuilder {
	q := github.Query("").WithRepo(repo).WithAuthor(auth.User).IsMerged()

	if !mergedAfter.IsZero() {
		log.Info("looking at pull requests merged after %s", mergedAfter.Format(commons.DateLayout))
		q.IsMergedAfter(mergedAfter)
	} else {
		log.Info("looking at pull requests by %s in %s", auth.User, repo)
	}

	return q
}

// runQuery searches for pull requests using a prepared query
func runQuery(client github.Client, repo string, query *github.QueryBuilder, page int) github.SearchResponse {
	log.Debug("getting page %d of pull requests", page)
	results, err := client.SearchPulls(repo, page, query)
	if err != nil {
		log.Fatal("failed to get pull requests: %s", err)
	}
	if results.IsIncomplete {
		log.Warn("search results are incomplete due to timeout")
	}

	return results
}

// summarize prints a summary of the collected stats
func summarize(repo string, count int, avgTimeToMerge time.Duration, avgReviews int, countByLabel map[github.Label]int) {
	log.Info("found %d merged pull requests for %s", count, repo)

	if count < 1 {
		return
	}

	avgTimeToMergeFriendly := durafmt.Parse(avgTimeToMerge)
	log.Info("average time to merge: %0.4f hours (%s)", avgTimeToMerge.Hours(), avgTimeToMergeFriendly)
	log.Info("average reviewers count: %d", avgReviews)
	for label, count := range countByLabel {
		log.Info("with label '%s': %d", label, count)
	}
}
