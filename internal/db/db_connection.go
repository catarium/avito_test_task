package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", "postgres://test:test@db:5432/test?sslmode=disable")
	return db, err
}
