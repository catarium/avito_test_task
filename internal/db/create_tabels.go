package db

import (
	"context"
	"database/sql"
)

func CreateTables(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS teams (
			team_name TEXT primary key
		);
		CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			user_name TEXT,
			team_name TEXT,
			is_active BOOLEAN,
			FOREIGN KEY (team_name) REFERENCES teams (team_name)
		);
		CREATE TABLE IF NOT EXISTS pull_requests (
			pull_request_id INT,
			pull_reqeust_name TEXT,
			author_id INT,
			is_merged BOOLEAN,
			created_at TIMESTAMP WITHOUT TIME ZONE,
			mereged_at TIMESTAMP WITHOUT TIME ZONE,
			FOREIGN KEY (author_id) REFERENCES users (user_id)
		);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}
