// internal/handlers/tasks_test.go
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/HywlEch/Todo_list/internal/models"
	"github.com/HywlEch/Todo_list/internal/store"
	"github.com/stretchr/testify/assert"
)

// TestGetTaskByID_Success 测试获取单个任务的“成功”路径
func TestGetTaskByID_Success(t *testing.T) {
	// --- ARRANGE (准备) ---
	// 1. 将 Gin 设置为测试模式
	gin.SetMode(gin.TestMode)

	// 2. 创建我们的模拟 store
	mockStore := new(store.MockStore)

	// 3. 定义预期的输入和输出
	mockTask := &models.Task{
		ID:        1,
		Title:     "Test Task",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// 4. “教” mockStore 如何行动：
	// 当 GetTaskByID 方法被以参数 `1` 调用时，返回 `mockTask` 并且不返回错误 (nil)
	mockStore.On("GetTaskByID", 1).Return(mockTask, nil)

	// 5. 用我们的 mock store 创建 handler
	taskHandler := NewTaskHandler(mockStore, nil)

	// --- ACT (执行) ---
	// 1. 设置路由
	router := gin.Default()
	router.GET("/tasks/:id", taskHandler.GetTaskByID)

	// 2. 创建一个假的 HTTP 请求
	req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	// 创建一个 ResponseRecorder 来记录响应
	w := httptest.NewRecorder()

	// 3. 让服务器处理这个假请求
	router.ServeHTTP(w, req)

	// --- ASSERT (断言) ---
	// 1. 断言 HTTP 状态码是 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. 断言响应体是我们期望的 JSON
	var responseTask models.Task
	err := json.Unmarshal(w.Body.Bytes(), &responseTask)
	assert.NoError(t, err) // 确认 JSON 解析没有出错
	assert.Equal(t, mockTask.ID, responseTask.ID)
	assert.Equal(t, mockTask.Title, responseTask.Title)
	
	// 3. 断言 mock store 的预期调用都已发生
	mockStore.AssertExpectations(t)
}

// TestGetTaskByID_NotFound 测试任务未找到的场景
func TestGetTaskByID_NotFound(t *testing.T) {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	mockStore := new(store.MockStore)

	mockStore.On("GetTaskByID", 2).Return(nil, store.ErrNotFound)
	taskHandler := NewTaskHandler(mockStore, nil)
	// ACT
	router := gin.Default()
	router.GET("/tasks/:id", taskHandler.GetTaskByID)
	req, _ := http.NewRequest(http.MethodGet, "/tasks/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ASSERT
	// 断言 HTTP 状态码是 404 Not Found
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockStore.AssertExpectations(t)
}