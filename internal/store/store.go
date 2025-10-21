package store

import(
	"context"
	"errors"
	"github.com/HywlEch/Todo_list/internal/models"
)

var ErrNotFound = errors.New("requested resource not found")
// Store 是我们数据存储层的接口
type Store interface {
	CreateTask(ctx context.Context,task *models.Task) error
	GetTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id int) error
}