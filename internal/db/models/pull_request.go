package models

import (
	"time"
)

type PullRequest struct {
	PullRequestId   int
	PullRequestName string
	AuthorId        int
	IsMerged        bool
	Reviewers       *[]User
	CreatedAt       time.Time
	MergedAt        time.Time
}
