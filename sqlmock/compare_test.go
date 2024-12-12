package sqlmock

import "testing"

func TestSQLCompare(t *testing.T) {
	// 语义相等的查询
	sql1 := "SELECT id, name FROM users WHERE age > 18 ORDER BY name"
	sql2 := "SELECT name, id FROM users WHERE age > 18 ORDER BY name"

	// 语义不等的查询
	sql3 := "SELECT id FROM users WHERE age > 21"

	// 比较
	if !CompareSQL(sql1, sql2) {
		t.Errorf("SQL statements should be semantically equal")
	}

	if CompareSQL(sql1, sql3) {
		t.Errorf("SQL statements should not be semantically equal")
	}
}
