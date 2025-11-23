package db

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS teams (
			team_name TEXT primary key
		);
		CREATE TABLE IF NOT EXISTS users (
			user_id TEXT PRIMARY KEY,
			username TEXT,
			team_name TEXT,
			is_active BOOLEAN,
			FOREIGN KEY (team_name) REFERENCES teams (team_name)
		);
		CREATE TABLE IF NOT EXISTS pull_requests (
			pull_request_id TEXT PRIMARY KEY,
			pull_request_name TEXT,
			author_id TEXT,
			is_merged BOOLEAN,
			created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
			merged_at TIMESTAMP WITHOUT TIME ZONE,
			FOREIGN KEY (author_id) REFERENCES users (user_id)
		);
		CREATE TABLE IF NOT EXISTS reviewers (
			pull_request_id TEXT,
			user_id TEXT,
			PRIMARY KEY (pull_request_id, user_id),
			FOREIGN KEY (pull_request_id) REFERENCES pull_requests (pull_request_id),
			FOREIGN KEY (user_id) REFERENCES users (user_id)
		);
	`
	_, err := db.Exec(query)
	return err
}
