package internal

import "database/sql"

// SqlDbFacade is an interface to allow functions to receive either a database connection or a database transaction.
type SqlDbFacade interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}
