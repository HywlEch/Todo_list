// internal/handlers/tasks.go
package handlers

import (
	//"database/sql"
	// "errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"fmt"
	

	"github.com/HywlEch/Todo_list/internal/apperrors"
	"github.com/HywlEch/Todo_list/internal/models"
	"github.com/HywlEch/Todo_list/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
)

//包含任务相关的 handler
type TaskHandler struct {
	Store 	store.Store 
	Redsync *redsync.Redsync
}

//创建一个新的 TaskHandler
func NewTaskHandler(s store.Store, rs *redsync.Redsync) *TaskHandler {
	return &TaskHandler{Store: s,
	Redsync: rs,
	}
}

//辅助函数 从Gin上下文中安全的获取userID
func getUserIDFromContext(c *gin.Context)(int, bool){
	userIDAny, ok := c.Get("userID")
	if !ok {
		c.Error(apperrors.NewUnauthorizedError("用户ID未找到", nil))
		return 0, false
	}
	userID, OK := userIDAny.(int)
	if !ok {
		c.Error(apperrors.NewInternalServerError("ID类型错误", nil))
		return 0, false
	}
	return userID, OK
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.Error(apperrors.NewBadRequestError("输入格式错误", err))
		return
	}
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	task.UserID = userID
	
	if err := h.Store.CreateTask(c.Request.Context(), &task); err != nil {
		c.Error(err)
		return
	}
	go h.sendNotification(&task)

	log.Printf("Created task with ID %d", task.ID)
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	tasks, err := h.Store.GetTasks(c.Request.Context(),userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(apperrors.NewBadRequestError("id格式错误", err))
		return
	}
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	task, err := h.Store.GetTaskByID(c.Request.Context(), id, userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(apperrors.NewBadRequestError("ID不合理", err))
		return
	}
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil { 
		c.Error(apperrors.NewBadRequestError("不合理得输入", err))
		return
	}
	task.ID = id
	task.UserID = userID

	//添加分布式锁
	lockKey := fmt.Sprintf("lock:task:%d", id)
	mutex := h.Redsync.NewMutex(lockKey, redsync.WithContext(c.Request.Context()))
	if err := mutex.Lock(); err != nil {
		log.Printf("获取锁失败: %v", err)
		c.Error(apperrors.NewInternalServerError("请稍后重试", err))
		return
	}
	log.Printf("获取锁成功")
	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil { 
			log.Printf("释放锁失败: %v", err)
		}else {
			log.Printf("释放锁成功")
		}
	}()

	if err := h.Store.UpdateTask(c.Request.Context(), &task); err != nil { 
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(apperrors.NewBadRequestError("ID格式错误", err))
		return
	}
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}
	if err := h.Store.DeleteTask(c.Request.Context(), id, userID); err != nil { 
		c.Error(err)
		return
	}
	//删除成功返回204
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) sendNotification(task *models.Task){
	log.Printf("Starting notification for task:%d(%s)", task.ID, task.Title)

	//模拟一个耗时操作
	time.Sleep(3 *time.Second)
	log.Printf("Successfully sent notification for task:%d(%s)", task.ID, task.Title)
	log.Printf("Successfully sent nitification for task:%d(%s)", task.ID, task.Title)
}