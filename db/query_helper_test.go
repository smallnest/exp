package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustDB(t testing.TB, schema string) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(schema)
	require.NoError(t, err)
	return db
}

func exampleDB(t *testing.T) *sql.DB {
	return mustDB(t, `CREATE TABLE persons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(120) NOT NULL DEFAULT ''
	);
	INSERT INTO PERSONS (name)
	VALUES ('brett'), ('fred');`)
}

type person struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `json:"name,omitempty"` // `db:"name" json:"name,omitempty"`
}

func TestRows(t *testing.T) {
	db := exampleDB(t)

	persons, err := Rows[person](db, "SELECT * FROM persons order by id")
	assert.NoError(t, err)
	require.Equal(t, 2, len(persons))
	assert.Equal(t, 1, persons[0].ID)
	assert.Equal(t, "brett", persons[0].Name)
	assert.Equal(t, 2, persons[1].ID)
	assert.Equal(t, "fred", persons[1].Name)

	names, err := Rows[string](db, "SELECT name FROM persons order by id")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(names))
	assert.Equal(t, "brett", names[0])
	assert.Equal(t, "fred", names[1])
}

func TestRow(t *testing.T) {
	db := exampleDB(t)

	person, err := Row[person](db, "SELECT * FROM persons order by id limit 1")
	assert.NoError(t, err)
	assert.Equal(t, 1, person.ID)
	assert.Equal(t, "brett", person.Name)

	name, err := Row[string](db, "SELECT name FROM persons order by id limit 1")
	assert.NoError(t, err)
	assert.Equal(t, "brett", name)
}
