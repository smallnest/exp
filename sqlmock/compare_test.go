package sqlmock

import "testing"

func TestSQLCompare(t *testing.T) {
	sql1 := "SELECT id, name FROM users WHERE age > 18 ORDER BY name"
	sql2 := "SELECT name, id FROM users WHERE age > 18 ORDER BY name"

	sql3 := "SELECT id FROM users WHERE age > 21"

	if !CompareSQL(sql1, sql2) {
		t.Errorf("SQL statements should be semantically equal")
	}

	if CompareSQL(sql1, sql3) {
		t.Errorf("SQL statements should not be semantically equal")
	}
}
