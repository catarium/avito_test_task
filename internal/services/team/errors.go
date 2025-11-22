package team

import (
	"fmt"

	"github.com/catarium/avito_test_task/internal/dto"
)

func ErrTeamExists(teamName string) *dto.ErrorDto {
	return &dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErrorTeamExists, Message: fmt.Sprintf("%s already exists", teamName)}}
}
