package sqlmock_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/smallnest/exp/sqlmock"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func (r *UserRepository) GetUserByID(id int) (User, error) {
	var user User
	err := r.db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Age)
	return user, err
}

func (r *UserRepository) CreateUser(user User) error {
	_, err := r.db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	return err
}

//----------------------------------------------------------

func TestUserRepository(t *testing.T) {
	// 1. 创建 mock 数据库
	mockDB := sqlmock.NewMock()

	// 2. 期望一个查询并设置返回值
	mockDB.ExpectQuery("SELECT id, name, age FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows([]string{"id", "name", "age"}, [][]driver.Value{
			{1, "John Doe", 30},
		})

	// 3. 打开数据库连接
	db, err := mockDB.Open("mock")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 4. 创建仓库实例
	repo := &UserRepository{db: db}

	// 5. 执行测试
	user, err := repo.GetUserByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 6. 验证结果
	if user.ID != 1 || user.Name != "John Doe" || user.Age != 30 {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestCreateUser(t *testing.T) {
	// 1. 创建 mock 数据库
	mockDB := sqlmock.NewMock()

	// 2. 期望一个插入操作
	mockDB.ExpectQuery("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 25).
		WillReturnRows([]string{"id"}, [][]driver.Value{
			{"2"}, // 返回插入的ID
		})

	// 3. 打开数据库连接
	db, err := mockDB.Open("mock")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 4. 创建仓库实例
	repo := &UserRepository{db: db}

	// 5. 执行测试
	err = repo.CreateUser(User{Name: "Alice", Age: 25})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestQueryError(t *testing.T) {
	// 1. 创建 mock 数据库
	mockDB := sqlmock.NewMock()

	// 2. 期望一个查询并返回错误
	mockDB.ExpectQuery("SELECT id, name, age FROM users WHERE id = ?", 999).
		WillReturnError([]string{"id", "name", "age"}, sql.ErrNoRows)

	// 3. 打开数据库连接
	db, err := mockDB.Open("mock")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 4. 创建仓库实例
	repo := &UserRepository{db: db}

	// 5. 执行测试
	_, err = repo.GetUserByID(999)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows, got %v", err)
	}
}