package pullrequest

import (
	"fmt"

	"github.com/catarium/avito_test_task/internal/dto"
)

func ErrPullRequestExists(pullRequestId string) *dto.ErrorDto {
	return &dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErrorPRExists, Message: fmt.Sprintf("%s already exists", pullRequestId)}}
}

var ErrPullRequestMerged dto.ErrorDto = dto.ErrorDto{Error: dto.ErrorDtoContent{
	Code:    dto.ErrorPRMerged,
	Message: "cannot reassign on merged PR",
}}

var ErrNotAssigned dto.ErrorDto = dto.ErrorDto{Error: dto.ErrorDtoContent{
	Code:    dto.ErrorNotAssigned,
	Message: "reviewer is not assigned to this PR",
}}

var ErrNoCandidate dto.ErrorDto = dto.ErrorDto{Error: dto.ErrorDtoContent{
	Code:    dto.ErrorNoCandidate,
	Message: "no active replacement candidate in team",
}}
