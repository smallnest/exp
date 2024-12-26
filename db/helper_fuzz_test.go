package db

import (
	"context"
	"database/sql"
	"testing"
)

// Fuzz test for Rows function
func FuzzRows(f *testing.F) {
	f.Add(1, "alice")

	f.Fuzz(func(t *testing.T, id int, name string) {
		db := exampleDB(t)
		defer db.Close()
		ctx := context.Background()

		// Create a dummy table and insert a row for testing
		_, err := db.ExecContext(ctx, "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
		if err != nil {
			t.Fatal(err)
		}

		Insert(ctx, db, "INSERT INTO users (id, name) VALUES (?, ?)", id, name)
		if err != nil {
			t.Fatal(err)
		}

		query := "SELECT * FROM users WHERE id = ? and name = ?"
		_, err = Rows[person](ctx, db, query, id, name)
		if err != nil && err != sql.ErrNoRows {
			t.Fatalf("Rows failed:id: %d, name: %s, err: %v", id, name, err)
		}

		count, err := Update(ctx, db, "UPDATE users SET name = ? WHERE id = ?", "bob", id)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("Update failed: id: %d, name: %s", id, name)
		}

		count, err = UpdateTx(ctx, db, "UPDATE users SET name = ? WHERE id = ?", "charlie", id)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("UpdateTx failed: id: %d, name: %s", id, name)
		}

		count, err = Delete(ctx, db, "DELETE FROM users WHERE id = ?", id)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("Delete failed: id: %d, name: %s", id, name)
		}

		query = "SELECT * FROM users WHERE id = ?"
		_, err = Row[person](ctx, db, query, id)
		if err == nil || err != sql.ErrNoRows {
			t.Fatalf("Rows failed:id: %d, name: %s, err: %v", id, name, err)
		}
	})
}
