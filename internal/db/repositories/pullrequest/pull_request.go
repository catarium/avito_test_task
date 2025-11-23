package pullrequest

import (
	"database/sql"

	"github.com/catarium/avito_test_task/internal/db/models"
	"github.com/catarium/avito_test_task/internal/utils/db"
)

type PullRequestRepository struct {
	DB *sql.DB
}

func (pr PullRequestRepository) toModel(row db.RowOrRows) (*models.PullRequest, error) {
	res := models.PullRequest{}
	err := row.Scan(&res.PullRequestId, &res.PullRequestName, &res.AuthorId, &res.IsMerged, &res.CreatedAt, &res.MergedAt)
	if err != nil {
		return nil, err
	}
	rows, err := pr.DB.Query("SELECT user_id FROM reviewers WHERE reviewers.pull_request_id = $1", res.PullRequestId)
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

func (pr PullRequestRepository) Create(pullRequestId string, pullRequestName string, authorId string) (*models.PullRequest, error) {
	query := `INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, is_merged)
	VALUES ($1, $2, $3, false)
	RETURNING pull_request_id, pull_request_name, author_id, is_merged, created_at, merged_at`
	row := pr.DB.QueryRow(query, pullRequestId, pullRequestName, authorId)
	return pr.toModel(row)
}

func (pr PullRequestRepository) Merge(pullRequestId string) (*models.PullRequest, error) {
	query := "UPDATE pull_requests SET is_merged = true, merged_at = now() WHERE pull_request_id = $1 RETURNING pull_request_id, pull_request_name, author_id, is_merged, created_at, merged_at"
	row := pr.DB.QueryRow(query, pullRequestId)
	return pr.toModel(row)
}

func (pr PullRequestRepository) Reassign(pullRequestId string, oldReviewerId string) (*models.PullRequest, string, error) {
	var new_user_id string
	query := `SELECT u1.user_id FROM users u1 JOIN users u2 ON u2.user_id = $1 JOIN pull_requests pr ON pr.pull_request_id = $2 WHERE
		u1.team_name = u2.team_name and
		u1.user_id != u2.user_id and
		u1.user_id != pr.author_id and
		u1.is_active = true and
		u1.user_id not in (SELECT user_id FROM reviewers WHERE pull_request_id = $2)`
	row := pr.DB.QueryRow(query, oldReviewerId, pullRequestId)
	err := row.Scan(&new_user_id)
	if err != nil {
		return nil, "", err
	}
	query = "UPDATE reviewers SET user_id = $1 WHERE pull_request_id = $2 and user_id = $3"
	_, err = pr.DB.Query(query, new_user_id, pullRequestId, oldReviewerId)
	if err != nil {
		return nil, "", err
	}
	query = "SELECT * FROM pull_requests WHERE pull_request_id = $1"
	row = pr.DB.QueryRow(query, pullRequestId)
	res, err := pr.toModel(row)
	if err != nil {
		return nil, "", err
	}
	return res, new_user_id, nil
}

func (pr PullRequestRepository) GetByReviewer(reviewerId string) ([]models.PullRequest, error) {
	res := []models.PullRequest{}
	query := `SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.is_merged
	FROM pull_requests pr JOIN reviewers r ON pr.pull_request_id = r.pull_request_id
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

func (pr PullRequestRepository) Exists(pullRequestId string) (bool, error) {
	query := "SELECT COUNT(*) FROM pull_requests WHERE pull_request_id = $1"
	row := pr.DB.QueryRow(query, pullRequestId)
	res := 1
	row.Scan(&res)
	if (row.Err() == nil) && (res > 0) {
		return true, nil
	} else if (row.Err() == sql.ErrNoRows) || (res == 0) {
		return false, nil
	}
	return false, row.Err()
}

func (pr PullRequestRepository) IsMerged(pullRequestId string) (bool, error) {
	query := "SELECT COUNT(*) FROM pull_requests WHERE pull_request_id = $1 AND is_merged = true"
	row := pr.DB.QueryRow(query, pullRequestId)
	res := 1
	row.Scan(&res)
	if (row.Err() == nil) && (res > 0) {
		return true, nil
	} else if (row.Err() == sql.ErrNoRows) || (res == 0) {
		return false, nil
	}
	return false, row.Err()
}

func (pr PullRequestRepository) IsAssigned(pullRequestId string, userId string) (bool, error) {
	query := "SELECT COUNT(*) FROM reviewers WHERE pull_request_id = $1 AND user_id = $2"
	row := pr.DB.QueryRow(query, pullRequestId, userId)
	res := 1
	row.Scan(&res)
	if (row.Err() == nil) && (res > 0) {
		return true, nil
	} else if (row.Err() == sql.ErrNoRows) || (res == 0) {
		return false, nil
	}
	return false, row.Err()
}

func (pr PullRequestRepository) Assign(pullRequestId string, n int) ([]string, error) {
	query := `INSERT INTO reviewers (pull_request_id, user_id) (
		SELECT $1::text, u.user_id FROM users u
			JOIN pull_requests pr ON pr.pull_request_id = $1::text
			JOIN users u2 ON u2.user_id = pr.author_id
			WHERE u.user_id != pr.author_id and u.is_active = true and u.team_name = u2.team_name and u.user_id not in
			(SELECT user_id FROM reviewers WHERE pull_request_id = $1::text)
		LIMIT $2
	) RETURNING user_id`
	rows, err := pr.DB.Query(query, pullRequestId, n)
	if err != nil {
		return nil, err
	}
	res := []string{}
	var el string
	for rows.Next() {
		err = rows.Scan(&el)
		if err != nil {
			return nil, err
		}
		res = append(res, el)
	}
	return res, nil
}
