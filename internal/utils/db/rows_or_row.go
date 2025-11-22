package db

type RowOrRows interface {
	Scan(dest ...any) error
}
