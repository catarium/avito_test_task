package dto

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            PRStatus `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"created_at"`
	MergedAt          string   `json:"merged_at"`
}
