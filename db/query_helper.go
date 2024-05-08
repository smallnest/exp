package db

import (
	"context"
	"database/sql"

	"github.com/blockloop/scan"
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

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return result, err
	}

	err = scan.Row(&result, rows)

	return result, err
}

// RowsMap is a helper function that wraps sql rows to scan into a slice of map[string]any.
// Key of the map is the column name, and value is the column value.
func RowsMap(ctx context.Context, db *sql.DB, query string, args ...any) ([]map[string]any, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]any, len(colNames))
	colPtrs := make([]any, len(colNames))

	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}

	var ret []map[string]any
	for rows.Next() {
		err = rows.Scan(colPtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range cols {
			row[colNames[i]] = col
		}
		ret = append(ret, row)
	}

	return ret, nil
}

// RowMap is a helper function that wraps sql rows to scan into a single map[string]any.
// Key of the map is the column name, and value is the column value.
// It returns the first row of the result set.
func RowMap(ctx context.Context, db *sql.DB, query string, args ...any) (map[string]any, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cols := make([]any, len(colNames))
	colPtrs := make([]any, len(colNames))

	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(colPtrs...)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]any)
	for i, col := range cols {
		ret[colNames[i]] = col
	}

	return ret, nil
}
