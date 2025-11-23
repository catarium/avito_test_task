package dto

type PullRequestReassign struct {
	Pr         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}
