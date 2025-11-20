package dto

type PullRequestShort struct {
	PullRequestId   int      `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorId        int      `json:"author_id"`
	Status          PRStatus `json:"status"`
}
