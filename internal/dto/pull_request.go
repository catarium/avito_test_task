package dto

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestId     int      `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          int      `json:"author_id"`
	Status            PRStatus `json:"status"`
	AssignedReviewers []int    `json:"assigned_reviewers"`
	CreatedAt         string   `json:"created_at"`
	MergedAt          string   `json:"merged_at"`
}
