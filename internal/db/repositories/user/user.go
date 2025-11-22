package user

import "database/sql"
import "github.com/catarium/avito_test_task/internal/db/models"

type UserRepository struct {
	DB *sql.DB
}

func (ur UserRepository) Create(userId string, username string, teamName string, isActive bool) (*models.User, error) {
	query := "INSERT INTO users (user_id, username, team_name, is_active) VALUES($1, $2, $3, $4)"
	row := ur.DB.QueryRow(query, userId, username, teamName, isActive)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return &models.User{UserId: userId, Username: username, TeamName: teamName, IsActive: isActive}, nil
}

func (ur UserRepository) Exists(userId string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE user_id = $1"
	row := ur.DB.QueryRow(query, userId)
	res := 1
	row.Scan(&res)
	if (row.Err() == nil) && (res > 0) {
		return true, nil
	} else if (row.Err() == sql.ErrNoRows) || (res == 0) {
		return false, nil
	}
	return false, row.Err()
}

func (ur UserRepository) Update(userId string, username string, teamName string, isActive bool) (*models.User, error) {
	query := "UPDATE users SET username=$2, team_name=$3, is_active=$4 WHERE user_id=$1;"
	row := ur.DB.QueryRow(query, userId, username, teamName, isActive)
	if row.Err() != nil {
		return nil, row.Err()
	}
	return &models.User{UserId: userId, Username: username, TeamName: teamName, IsActive: isActive}, nil
}

func (ur UserRepository) SetIsActive(userId string, isActive bool) (*models.User, error) {
	query := "UPDATE users SET is_active=$2 WHERE user_id=$1 RETURNING user_id, username, team_name, is_active"
	row := ur.DB.QueryRow(query, userId, isActive)
	if row.Err() != nil {
		return nil, row.Err()
	}
	res := models.User{}
	row.Scan(&res.UserId, &res.Username, &res.TeamName, &res.IsActive)
	return &res, nil
}

func (ur UserRepository) GetByTeamName(teamName string) ([]models.User, error) {
	res := []models.User{}
	query := "SELECT user_id, username, team_name, is_active FROM users WHERE team_name = $1"
	rows, err := ur.DB.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	var user models.User
	for rows.Next() {
		user = models.User{}
		err = rows.Scan(&user.UserId, &user.Username, &user.TeamName, &user.IsActive)
		res = append(res, user)
		if err != nil {
			return nil, err
		}
	}
	return res, err
}
