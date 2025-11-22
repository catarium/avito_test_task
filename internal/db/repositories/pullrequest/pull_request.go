package pullrequest

import (
	"context"
	"database/sql"

	"github.com/catarium/avito_test_task/internal/db/models"
	"github.com/catarium/avito_test_task/internal/utils/db"
)

type PullRequestRepository struct {
	DB *sql.DB
}

func (pr PullRequestRepository) toModel(row db.RowOrRows) (*models.PullRequest, error) {
	res := models.PullRequest{}
	err := row.Scan(&res.PullRequestId, &res.PullRequestName, &res.AuthorId, &res.IsMerged)
	if err != nil {
		return nil, err
	}
	rows, err := pr.DB.Query("SELECT user_id FROM reviews WHERE reviews.pull_request_id = $1", res.PullRequestId)
	if err != nil {
		return nil, err
	}
	var curr_reviewer string
	for rows.Next() {
		rows.Scan(&curr_reviewer)
		res.Reviewers = append(res.Reviewers, curr_reviewer)
	}
	return &res, nil
}

func (pr PullRequestRepository) Create(ctx context.Context, pullRequestId string, pullRequestName string, authorId string) (*models.PullRequest, error) {
	query := `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, is_merged)
	VALUES ($1, $2, $3, false)
	RETURNING pull_request_id, pull_request_name, author_id, is_merged, created_at, merged_at`
	row := pr.DB.QueryRow(query, pullRequestId, pullRequestName, authorId)
	return pr.toModel(row)
}

func (pr PullRequestRepository) Merge(pullRequestId string) (*models.PullRequest, error) {
	query := "UPDATE pull_requests SET is_merged = true WHERE pull_request_id = $1 RETURNING pull_request_id, pull_request_name, author_id, is_merged, created_at, merged_at"
	row := pr.DB.QueryRow(query, pullRequestId)
	return pr.toModel(row)
}

func (pr PullRequestRepository) Reassign(pullRequestId string, oldReviewerId string) (*models.PullRequest, error) {
	var new_user_id string
	query := `SELECT user_id FROM users u1 JOIN users u2 ON u2.user_id = $1 JOIN pull_requests pr ON pr.pull_request_id = $2 WHERE
		u1.team_name = u2.team_name and
		u1.user_id != u2.user_id and
		u1.user_id != pr.author_id
		u1.is_active = true`
	row := pr.DB.QueryRow(query, oldReviewerId, pullRequestId)
	err := row.Scan(&new_user_id)
	if err != nil {
		return nil, err
	}
	query = "UPDATE reviwers SET user_id = $1 WHERE pull_request_id = $2 and user_id = $3"
	_, err = pr.DB.Query(query, new_user_id, pullRequestId, oldReviewerId)
	if err != nil {
		return nil, err
	}
	query = "SELECT * FROM pull_requests WHERE pull_request_id = $1"
	row = pr.DB.QueryRow(query, pullRequestId)
	return pr.toModel(row)
}

func (pr PullRequestRepository) GetByReviewer(reviewerId string) ([]models.PullRequest, error) {
	res := []models.PullRequest{}
	query := `SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
	FROM pull_requests pr JOIN reviews r ON pr.pull_request_id = r.pull_request_id
	WHERE r.user_id = $1;`
	row, err := pr.DB.Query(query, reviewerId)
	if err != nil {
		return nil, err
	}
	var pullRequest models.PullRequest
	for row.Next() {
		row.Scan(&pullRequest.PullRequestId, &pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.IsMerged)
		res = append(res, pullRequest)
	}
	return res, nil
}
