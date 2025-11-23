package pullrequest

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/catarium/avito_test_task/internal/db/models"
	"github.com/catarium/avito_test_task/internal/db/repositories/pullrequest"
	"github.com/catarium/avito_test_task/internal/db/repositories/user"
	"github.com/catarium/avito_test_task/internal/dto"
	"github.com/catarium/avito_test_task/internal/services"
)

type PullRequestService struct {
	PullRequestRepository *pullrequest.PullRequestRepository
	UserRepository        *user.UserRepository
}

func (ps PullRequestService) Create(pullRequestId string, pullRequestName string, authorId string) (*dto.PullRequestDto, *dto.ErrorDto, int) {
	exists, err := ps.UserRepository.Exists(authorId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	exists, err = ps.PullRequestRepository.Exists(pullRequestId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if exists {
		return nil, ErrPullRequestExists(pullRequestId), http.StatusConflict
	}
	pullRequest, err := ps.PullRequestRepository.Create(pullRequestId, pullRequestName, authorId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	res := dto.PullRequestDto{
		Pr: dto.PullRequestContentDto{
			PullRequestId:   pullRequest.PullRequestId,
			PullRequestName: pullRequest.PullRequestName,
			AuthorId:        pullRequest.AuthorId,
			CreatedAt:       pullRequest.CreatedAt.Time.Format(time.RFC3339),
			MergedAt:        pullRequest.MergedAt.Time.Format(time.RFC3339),
		},
	}
	if pullRequest.IsMerged {
		res.Pr.Status = dto.StatusMerged
	} else {
		res.Pr.Status = dto.StatusOpen
	}
	reviewers, err := ps.PullRequestRepository.Assign(pullRequestId, 2)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	res.Pr.AssignedReviewers = reviewers
	return &res, nil, http.StatusCreated
}

func (ps PullRequestService) Merge(pullRequestId string) (*dto.PullRequestDto, *dto.ErrorDto, int) {
	exists, err := ps.PullRequestRepository.Exists(pullRequestId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	merged, err := ps.PullRequestRepository.IsMerged(pullRequestId)
	var pullRequest *models.PullRequest
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if merged {
		pullRequest, err = ps.PullRequestRepository.GetByPullRequestId(pullRequestId)
	} else {
		pullRequest, err = ps.PullRequestRepository.Merge(pullRequestId)
	}
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	res := dto.PullRequestDto{
		Pr: dto.PullRequestContentDto{
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			AuthorId:          pullRequest.AuthorId,
			CreatedAt:         pullRequest.CreatedAt.Time.Format(time.RFC3339),
			MergedAt:          pullRequest.MergedAt.Time.Format(time.RFC3339),
			AssignedReviewers: pullRequest.Reviewers,
		},
	}
	if pullRequest.IsMerged {
		res.Pr.Status = dto.StatusMerged
	} else {
		res.Pr.Status = dto.StatusOpen
	}
	return &res, nil, http.StatusOK
}

func (ps PullRequestService) Reassign(pullRequestId string, oldReviewerId string) (*dto.PullRequestReassign, *dto.ErrorDto, int) {
	exists, err := ps.PullRequestRepository.Exists(pullRequestId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	merged, err := ps.PullRequestRepository.IsMerged(pullRequestId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if merged {
		return nil, &ErrPullRequestMerged, http.StatusConflict
	}
	exists, err = ps.UserRepository.Exists(oldReviewerId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	assigned, err := ps.PullRequestRepository.IsAssigned(pullRequestId, oldReviewerId)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !assigned {
		return nil, &ErrNotAssigned, http.StatusConflict
	}
	pullRequest, new_user_id, err := ps.PullRequestRepository.Reassign(pullRequestId, oldReviewerId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &ErrNoCandidate, http.StatusConflict
	}
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	res := dto.PullRequestReassign{
		Pr: dto.PullRequestContentDto{
			PullRequestId:     pullRequest.PullRequestId,
			PullRequestName:   pullRequest.PullRequestName,
			AuthorId:          pullRequest.AuthorId,
			CreatedAt:         pullRequest.CreatedAt.Time.Format(time.RFC3339),
			MergedAt:          pullRequest.MergedAt.Time.Format(time.RFC3339),
			AssignedReviewers: pullRequest.Reviewers,
		},
		ReplacedBy: new_user_id,
	}
	if pullRequest.IsMerged {
		res.Pr.Status = dto.StatusMerged
	} else {
		res.Pr.Status = dto.StatusOpen
	}
	return &res, nil, http.StatusOK
}
