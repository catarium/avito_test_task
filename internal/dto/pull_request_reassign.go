package dto

type PullRequestReassign struct {
	Pr         PullRequestContentDto `json:"pr"`
	ReplacedBy string                `json:"replaced_by"`
}
