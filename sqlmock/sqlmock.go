package sqlmock

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
)

var (
	_ driver.Connector = &MockConnector{}
	_ driver.Driver    = &MockDriver{}
	_ driver.Conn      = &MockConn{}
	_ driver.Stmt      = &MockStmt{}
	_ driver.Result    = &MockResult{}
	_ driver.Rows      = &MockRows{}
	_ driver.Tx        = &MockTx{}
)

// ExpectedQuery 表示一个预期的查询
type ExpectedQuery struct {
	query   string
	args    []driver.Value
	rows    [][]driver.Value
	columns []string
	err     error
}

// MockDB 模拟数据库连接
type MockDB struct {
	mu       sync.Mutex
	expected []*ExpectedQuery
}

// NewMock 创建一个新的模拟数据库
func NewMock() *MockDB {
	return &MockDB{
		expected: []*ExpectedQuery{},
	}
}

// ExpectQuery 期望一个特定的查询
func (m *MockDB) ExpectQuery(query string, args ...driver.Value) *ExpectedQuery {
	m.mu.Lock()
	defer m.mu.Unlock()

	eq := &ExpectedQuery{
		query: query,
		args:  args,
	}
	m.expected = append(m.expected, eq)
	return eq
}

func (eq *ExpectedQuery) WithArgs(args ...driver.Value) *ExpectedQuery {
	eq.args = args
	return eq
}

// WillReturnRows 为查询设置返回的行数据
func (eq *ExpectedQuery) WillReturnRows(columns []string, rows [][]driver.Value) *ExpectedQuery {
	eq.rows = rows
	eq.columns = columns
	return eq
}

// WillReturnError 为查询设置返回的错误
func (eq *ExpectedQuery) WillReturnError(columns []string, err error) *ExpectedQuery {
	eq.err = err
	eq.columns = columns
	return eq
}

// Open 模拟数据库连接
func (m *MockDB) Open(driverName string) (*sql.DB, error) {
	connector := &MockConnector{mockDB: m}
	return sql.OpenDB(connector), nil
}

// MockConnector 实现 driver.Connector 接口
type MockConnector struct {
	mockDB *MockDB
}

func (mc *MockConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return &MockConn{mockDB: mc.mockDB}, nil
}

func (mc *MockConnector) Driver() driver.Driver {
	return &MockDriver{mockDB: mc.mockDB}
}

// MockDriver 实现 driver.Driver 接口
type MockDriver struct {
	mockDB *MockDB
}

func (md *MockDriver) Open(name string) (driver.Conn, error) {
	return &MockConn{mockDB: md.mockDB}, nil
}

// MockConn 实现 driver.Conn 接口
type MockConn struct {
	mockDB *MockDB
}

func (mc *MockConn) Prepare(query string) (driver.Stmt, error) {
	return &MockStmt{
		mockDB: mc.mockDB,
		query:  query,
	}, nil
}

func (mc *MockConn) Close() error {
	return nil
}

func (mc *MockConn) Begin() (driver.Tx, error) {
	return &MockTx{}, nil
}

// MockStmt 实现 driver.Stmt 接口
type MockStmt struct {
	mockDB *MockDB
	query  string
}

func (ms *MockStmt) Close() error {
	return nil
}

func (ms *MockStmt) NumInput() int {
	return -1
}

func (ms *MockStmt) Exec(args []driver.Value) (driver.Result, error) {
	ms.mockDB.mu.Lock()
	defer ms.mockDB.mu.Unlock()

	for i, expected := range ms.mockDB.expected {
		if CompareSQL(expected.query, ms.query) && matchArgs(expected.args, args) {
			ms.mockDB.expected = append(ms.mockDB.expected[:i], ms.mockDB.expected[i+1:]...)

			if expected.err != nil {
				return nil, expected.err
			}

			return &MockResult{}, nil
		}
	}

	return nil, fmt.Errorf("unexpected query: %s", ms.query)
}

func (ms *MockStmt) Query(args []driver.Value) (driver.Rows, error) {
	ms.mockDB.mu.Lock()
	defer ms.mockDB.mu.Unlock()

	for i, expected := range ms.mockDB.expected {
		if CompareSQL(expected.query, ms.query) && matchArgs(expected.args, args) {
			ms.mockDB.expected = append(ms.mockDB.expected[:i], ms.mockDB.expected[i+1:]...)

			if expected.err != nil {
				return nil, expected.err
			}

			return &MockRows{columns: expected.columns, rows: expected.rows}, nil
		}
	}

	return nil, fmt.Errorf("unexpected query: %s", ms.query)
}

// MockResult 实现 driver.Result 接口
type MockResult struct{}

func (mr *MockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (mr *MockResult) RowsAffected() (int64, error) {
	return 0, nil
}

// MockRows 实现 driver.Rows 接口
type MockRows struct {
	rows    [][]driver.Value
	columns []string
	cursor  int
}

func (mr *MockRows) Columns() []string {
	return mr.columns
}

func (mr *MockRows) Close() error {
	return nil
}

func (mr *MockRows) Next(dest []driver.Value) error {
	if mr.cursor >= len(mr.rows) {
		return sql.ErrNoRows
	}

	copy(dest, mr.rows[mr.cursor])
	mr.cursor++
	return nil
}

// MockTx 实现 driver.Tx 接口
type MockTx struct{}

func (mt *MockTx) Commit() error {
	return nil
}

func (mt *MockTx) Rollback() error {
	return nil
}

// 辅助函数：匹配查询参数
func matchArgs(expected, actual []driver.Value) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		a, _ := json.Marshal(expected[i])
		b, _ := json.Marshal(actual[i])

		if !bytes.Equal(a, b) {
			return false
		}
	}

	return true
}
