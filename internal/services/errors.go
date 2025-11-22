package services

import "github.com/catarium/avito_test_task/internal/dto"

func ErrUnknown(msg string) *dto.ErrorDto {
	return &dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErroUnknown, Message: msg}}
}

var ErrNotFound = dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErrorNotFound, Message: "resource not found"}}

func ErrInvalidJson(msg string) *dto.ErrorDto {
	return &dto.ErrorDto{Error: dto.ErrorDtoContent{Code: dto.ErrorInvalidBody, Message: msg}}
}
