package db

import (
	"database/sql"

	"github.com/blockloop/scan"
)

// Rows is a helper function that wraps sql rows to scan into a slice.
// Rows scans structs based on their db tag, and scan any fields not tagged with the db tag matched column name,
// So had better to tag all fields with db tag to avoid unexpected behavior.
func Rows[T any](db *sql.DB, query string, args ...any) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var result []T
	err = scan.Rows(&result, rows)

	return result, err
}

// Row is a helper function that wraps sql rows to scan into a single struct.
// Row scans structs based on their db tag, and scan any fields not tagged with the db tag matched column name,
// So had better to tag all fields with db tag to avoid unexpected behavior.
func Row[T any](db *sql.DB, query string, args ...any) (T, error) {
	var result T

	rows, err := db.Query(query, args...)
	if err != nil {
		return result, err
	}

	err = scan.Row(&result, rows)

	return result, err
}
