package models

import (
	"time"
)

type PullRequest struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	IsMerged        bool
	Reviewers       []string
	CreatedAt       time.Time
	MergedAt        time.Time
}
