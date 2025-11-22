package dto

type PullRequestShortDto struct {
	PullRequestId   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorId        string   `json:"author_id"`
	Status          PRStatus `json:"status"`
}

type UserPullRequestShortDto struct {
	UserId       string                `json:"user_id"`
	PullRequests []PullRequestShortDto `json:"pull_requests"`
}
