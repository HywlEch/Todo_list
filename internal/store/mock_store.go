package store

import (
	"context"
	"github.com/HywlEch/Todo_list/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock //嵌入 testify 的mock对象
}

//模拟CreateTask实现
func (m *MockStore) CreateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

// GetTasks 的模拟实现
func (m *MockStore) GetTasks(ctx context.Context) ([]models.Task, error){
	args := m.Called(ctx)
	return args.Get(0).([]models.Task), args.Error(1)
}

// GetTaskByID 的模拟实现
func (m *MockStore) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

// UpdateTask 是模拟实现
func (m *MockStore) UpdateTask(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

// DeleteTask 是模拟实现
func (m *MockStore) DeleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}