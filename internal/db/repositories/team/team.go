package team

import (
	"database/sql"

	"github.com/catarium/avito_test_task/internal/db/models"
)

type TeamRepository struct {
	DB *sql.DB
}

func (tr TeamRepository) CreateTeam(teamName string) (*models.Team, error) {
	query := "INSERT INTO teams (team_name) VALUES ($1) RETURNING team_name;"
	_, err := tr.DB.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	return &models.Team{TeamName: teamName}, nil
}

func (tr TeamRepository) Exists(teamName string) (bool, error) {
	query := "SELECT COUNT(*) FROM teams WHERE team_name = $1;"
	row := tr.DB.QueryRow(query, teamName)
	res := 1
	row.Scan(&res)
	if (row.Err() == nil) && (res > 0) {
		return true, nil
	} else if (row.Err() == sql.ErrNoRows) || (res == 0) {
		return false, nil
	}
	return false, row.Err()
}
