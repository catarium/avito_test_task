package schemas

type PullRequestReassignSchema struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}
