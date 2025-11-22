package team

import (
	"fmt"

	"github.com/catarium/avito_test_task/internal/dto"
)

var ErrTeamNotFound = dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErrorNotFound, Message: "resource not found"}}

func ErrUnknown(msg string) *dto.ErrorDto {
	return &dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErroUnknown, Message: msg}}
}

func ErrTeamExists(teamName string) *dto.ErrorDto {
	return &dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErrorTeamExists, Message: fmt.Sprintf("%s already exists", teamName)}}
}
