// internal/handlers/tasks.go
package handlers

import (
	//"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/HywlEch/Todo_list/internal/models"
	"github.com/HywlEch/Todo_list/internal/store"
	"github.com/gin-gonic/gin"
)

// TaskHandler 包含任务相关的 handler
type TaskHandler struct {
	Store store.Store // 依赖接口，而不是具体实现
}

// NewTaskHandler 创建一个新的 TaskHandler
func NewTaskHandler(s store.Store) *TaskHandler {
	return &TaskHandler{Store: s}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Store.CreateTask(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	go h.sendNotification(&task)

	log.Printf("Created task with ID %d", task.ID)
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.Store.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	task, err := h.Store.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	var task models.Task
	task.ID = id
	if err := c.ShouldBindJSON(&task); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Store.UpdateTask(c.Request.Context(), &task); err != nil { 
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	if err := h.Store.DeleteTask(c.Request.Context(), id); err != nil { 
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
	}
}

func (h *TaskHandler) sendNotification(task *models.Task){
	log.Printf("Starting notification for task:%d(%s)", task.ID, task.Title)

	//模拟一个耗时操作
	time.Sleep(3 *time.Second)

	log.Printf("Successfully sent nitification for task:%d(%s)", task.ID, task.Title)
}