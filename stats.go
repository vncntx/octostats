package main

import (
	"fmt"
	"time"

	"octostats/commons"
	"octostats/github"

	"github.com/hako/durafmt"
)

// stats are information about a pull request
type stats struct {
	number      int            // pull request identifier
	mergedAt    time.Time      // time when the pull request was merged
	timeToMerge time.Duration  // time to merge from creation
	reviews     int            // number of reviews
	labels      []github.Label // labels
}

// log the stats
func (s stats) log() {
	timeToMergeFriendly := durafmt.Parse(s.timeToMerge).LimitFirstN(2)
	log.Info("pull request #%d took %s to merge on %s (%d reviews)", s.number, timeToMergeFriendly, s.mergedAt.Format(commons.DateLayout), s.reviews)
}

// getStats returns stats about a pull request search result
func getStats(client github.Client, repo string, item github.SearchResult) (stats, error) {
	result := stats{}

	// get pull request details
	pull, err := client.GetPull(repo, item.Number)
	if err != nil {
		return result, fmt.Errorf("couldn't get details of pull request #%d: %w", item.Number, err)
	}

	result.number = pull.Number
	result.mergedAt = pull.MergedAt

	// get time to merge
	if pull.CreatedAt.IsZero() || pull.MergedAt.IsZero() {
		return result, fmt.Errorf("pull request #%d has invalid timestamps", pull.Number)
	}
	result.timeToMerge = pull.MergedAt.Sub(pull.CreatedAt)

	// get number of reviews
	reviews, err := countReviews(client, repo, pull.Number)
	if err != nil {
		return result, fmt.Errorf("couldn't count reviews for pull request #%d: %w", pull.Number, err)
	}

	result.reviews = reviews
	result.labels = pull.Labels

	return result, nil
}

// countReviews counts the number of reviews for a pull request
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
