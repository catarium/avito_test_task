package models

import (
	"database/sql"
)

type PullRequest struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	IsMerged        bool
	Reviewers       []string
	CreatedAt       sql.NullTime
	MergedAt        sql.NullTime
}
