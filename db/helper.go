package db

import (
	"context"
	"database/sql"
	"strings"

	"github.com/blockloop/scan/v2"
)

// Rows is a helper function that wraps sql rows to scan into a slice.
// Rows scans structs based on their db tag, and scan any fields not tagged with the db tag matched column name,
// So had better to tag all fields with db tag to avoid unexpected behavior.
func Rows[T any](ctx context.Context, db *sql.DB, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
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
func Row[T any](ctx context.Context, db *sql.DB, query string, args ...any) (T, error) {
	var result T

	query = strings.TrimSuffix(strings.TrimSpace(query), ";")

	if !strings.Contains(strings.ToUpper(query), "LIMIT") {
		query += " LIMIT 1"
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return result, err
	}

	err = scan.Row(&result, rows)

	return result, err
}

// Count is a helper function that wraps sql rows to scan into a single int.
// The query should return a single column with interger type such as count,sum etc. with a single row.
func Count(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	// it is a query more simple than Row[int64]
	var count int64

	err := db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Insert is a helper function that wraps sql exec to insert a row.
func Insert(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// InsertTx is a helper function that wraps sql exec to insert a row in a transaction.
func InsertTx(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Delete is a helper function that wraps sql exec to delete rows.
func Delete(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Update is a helper function that wraps sql exec to update rows.
func Update(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// UpdateTx is a helper function that wraps sql exec to update rows in a transaction.
func UpdateTx(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
