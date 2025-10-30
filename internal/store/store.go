package store

import(
	"context"
	"errors"
	"github.com/HywlEch/Todo_list/internal/models"
)

var ErrNotFound = errors.New("requested resource not found")
var ErrUserExists = errors.New("user already exists")
// Store 是我们数据存储层的接口
type Store interface {
	CreateUser(ctx context.Context,user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	CreateTask(ctx context.Context,task *models.Task) error
	GetTasks(ctx context.Context, userId int) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id int, userId int) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id int, userId int) error
}