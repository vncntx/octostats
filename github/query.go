package github

import (
	"fmt"
	"strings"
	"time"

	"vincent.click/pkg/octostats/commons"
)

// QueryBuilder creates a search query with keywords and qualifiers
// See https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests
type QueryBuilder struct {
	keywords    string
	typ         string
	state       string
	author      string
	repo        string
	isDraft     bool
	isMerged    bool
	mergedAfter time.Time
}

// Query returns a new QueryBuilder
func Query(keywords string) *QueryBuilder {
	return &QueryBuilder{
		keywords: keywords,
	}
}

// IsDraft adds a draft pull request qualifier
func (q *QueryBuilder) IsDraft() *QueryBuilder {
	q.isDraft = true

	return q
}

// IsMerged adds a merged pull request qualifier
func (q *QueryBuilder) IsMerged() *QueryBuilder {
	q.isMerged = true

	return q
}

// WithState filters issues and pull requests by state
func (q *QueryBuilder) WithState(state string) *QueryBuilder {
	q.state = state

	return q
}

// WithAuthor filters issues and pull requests by the author's login username
func (q *QueryBuilder) WithAuthor(authorLogin string) *QueryBuilder {
	q.author = authorLogin

	return q
}

// WithRepo filters issues and pull requests by repository
func (q *QueryBuilder) WithRepo(repo string) *QueryBuilder {
	q.repo = repo

	return q
}

// IsMergedAfter filters pull requests that were merged after a given date
func (q *QueryBuilder) IsMergedAfter(t time.Time) *QueryBuilder {
	q.mergedAfter = t

	return q
}

// withType adds a type qualifier
func (q *QueryBuilder) withType(typ string) *QueryBuilder {
	q.typ = typ

	return q
}

// Build returns the query as a string
func (q *QueryBuilder) Build() string {
	parts := []string{}

	if len(q.typ) > 0 {
		parts = append(parts, is(q.typ))
	}
	if q.isDraft {
		parts = append(parts, is(draft))
	}
	if q.isMerged {
		parts = append(parts, is(merged))
	}
	if len(q.state) > 0 {
		parts = append(parts, is(q.state))
	}
	if len(q.author) > 0 {
		parts = append(parts, has("author", q.author))
	}
	if len(q.repo) > 0 {
		parts = append(parts, has("repo", q.repo))
	}
	if !q.mergedAfter.IsZero() {
		date := q.mergedAfter.Format(commons.DateLayout)
		parts = append(parts, fmt.Sprintf("merged:>%s", date))
	}
	if len(q.keywords) > 0 {
		parts = append(parts, q.keywords)
	}

	return strings.Join(parts, " ")
}

func is(adjective string) string {
	return fmt.Sprintf("is:%s", adjective)
}

func has(prop, value string) string {
	return fmt.Sprintf("%s:%s", prop, value)
}
