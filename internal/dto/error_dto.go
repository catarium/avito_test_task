package dto

type ErrorCode string

const (
	ErrorTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrorPRExists    ErrorCode = "PR_EXISTS"
	ErrorPRMerged    ErrorCode = "PR_MERGED"
	ErrorNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrorNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrorNotFound    ErrorCode = "NOT_FOUND"
	ErrorInvalidBody ErrorCode = "INVALID_BODY"
	ErroUnknown      ErrorCode = "UNKNOWN"
)

type ErrorDtoContent struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorDto struct {
	Error ErrorDtoContent `json:"error"`
}
