package user

import (
	"net/http"

	"github.com/catarium/avito_test_task/internal/db/repositories/pullrequest"
	"github.com/catarium/avito_test_task/internal/db/repositories/user"
	"github.com/catarium/avito_test_task/internal/dto"
	"github.com/catarium/avito_test_task/internal/services"
)

type UserService struct {
	UserRepository        user.UserRepository
	PullRequestRepository pullrequest.PullRequestRepository
}

func (us UserService) SetActive(userId string, isActive bool) (*dto.UserDto, *dto.ErrorDto, int) {
	exists, err := us.UserRepository.Exists(userId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	user, err := us.UserRepository.SetIsActive(userId, isActive)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	return &dto.UserDto{
			User: dto.UserDtoContent{
				UserId:   user.UserId,
				Username: user.Username,
				TeamName: user.TeamName,
				IsActive: user.IsActive,
			}},
		nil,
		http.StatusOK
}

func (us UserService) GetReviewedPullRequestsByUserId(userId string) (*dto.UserPullRequestShortDto, *dto.ErrorDto, int) {
	res := dto.UserPullRequestShortDto{UserId: userId}
	exists, err := us.UserRepository.Exists(userId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	pullRequests, err := us.PullRequestRepository.GetByReviewer(userId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	pullRequest := dto.PullRequestShortDto{}
	for _, pr := range pullRequests {
		pullRequest.PullRequestId = pr.PullRequestId
		pullRequest.PullRequestName = pr.PullRequestName
		pullRequest.AuthorId = pr.AuthorId
		if pr.IsMerged {
			pullRequest.Status = dto.StatusMerged
		} else {
			pullRequest.Status = dto.StatusOpen
		}
		res.PullRequests = append(res.PullRequests, pullRequest)
	}
	return &res, nil, http.StatusOK
}
