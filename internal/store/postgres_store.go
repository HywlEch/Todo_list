// internal/store/postgres_store.go
package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/HywlEch/Todo_list/internal/config"
	"github.com/HywlEch/Todo_list/internal/models"
	_ "github.com/lib/pq"
)

// PostgresStore 实现了 Store 接口
type PostgresStore struct {
	DB *sql.DB
}

// NewPostgresStore 创建一个新的 PostgresStore 实例
func NewPostgresStore(cfg config.DBConfig) (*PostgresStore, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.User, cfg.Password, cfg.DBName, cfg.Host, cfg.Port, cfg.SSLMode)
	// connStr := "user=todouser password=todopass dbname=todolist_db host=localhost port=5433 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{DB: db}, nil
}

func (s *PostgresStore) CreateTask(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (title, content, done) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;`
	err:= s.DB.QueryRowContext(ctx, query, task.Title, task.Content, task.Done).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil { 
		return fmt.Errorf("创建任务失败: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetTasks(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, title, content, done, created_at, updated_at FROM tasks ORDER BY created_at DESC;`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("store: failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.Done, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: failed to scan task row: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *PostgresStore) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	query := `SELECT id, title, content, done, created_at, updated_at FROM tasks WHERE id = $1;`
	var task models.Task
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.Content, &task.Done, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound 
		}
		return nil, fmt.Errorf("store: failed to update task %d: %w", task.ID, err)
	}
	return &task, nil
}

func (s *PostgresStore) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET title = $1, content = $2, done = $3 WHERE id = $4 RETURNING updated_at;`
	// 我们需要扫描返回的 updated_at，更新到传入的 task 对象上
	err := s.DB.QueryRowContext(ctx, query, task.Title, task.Content, task.Done, task.ID).Scan(&task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1;`
	_, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("store: failed to delete task %d: %w", id, err)
	}
	return nil
}