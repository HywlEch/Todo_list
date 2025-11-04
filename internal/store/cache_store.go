package store

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"context"

	"github.com/HywlEch/Todo_list/internal/models"
	"github.com/go-redis/redis/v8"
)

type CacheStore struct {
	redisClient *redis.Client
	next        Store
	ttl         time.Duration
}

// 创建一个新的CacheStore实例
func NewCacheStore(nextStore Store, redis *redis.Client) *CacheStore {
	return &CacheStore{
		next:        nextStore,
		redisClient: redis,
		ttl:         1 * time.Hour}
}

// 键名辅助函数
func taskKey(id int) string {
	return fmt.Sprintf("task:%d", id)
}

func userTaskKey(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}

// 缓存核心逻辑
func (s *CacheStore) GetTaskByID(ctx context.Context, id int, userID int) (*models.Task, error) {
	key := taskKey(id)

	//读缓存，尝试从redis中获取
	val, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		var task models.Task
		if err := json.Unmarshal([]byte(val), &task); err == nil {
			if task.UserID == userID {
				return &task, nil
			}
		}
	}
	if err != redis.Nil {
		//如果错误不是“key not found”Redis出错
		log.Printf("[CacheStore]Warn:Redis Fet error on key%s: %v", key, err)
	}
	log.Printf("[CacheStore]MISS: GetTaskByID(key: %s)", key)

	//调用真正的数据库
	task, err := s.next.GetTaskByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(task)
	if err != nil {
		log.Printf("[CacheStore]Error: Failed to marshal task: %v", err)
		return task, err
	}
	if err := s.redisClient.Set(ctx, key, jsonData, s.ttl).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to set task in redis: %v", err)
	}
	return task, nil
}

// 缓存核心逻辑
func (s *CacheStore) GetTasks(ctx context.Context, userID int) ([]models.Task, error) {
	key := userTaskKey(userID)
	val, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		var tasks []models.Task
		if err := json.Unmarshal([]byte(val), &tasks); err == nil {
			return tasks, nil
		}
	}
	log.Printf("[CacheStore]MISS: GetTasks(key: %s)", key)
	tasks, err := s.next.GetTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("[CacheStore]Error: Failed to marshal tasks: %v", err)
		return tasks, err
	}
	if err := s.redisClient.Set(ctx, key, jsonData, s.ttl).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to set tasks in redis: %v", err)
		return tasks, err
	}
	return tasks, nil
}

// 缓存失效逻辑
func (s *CacheStore) CreateTask(ctx context.Context, task *models.Task) error {
	err := s.next.CreateTask(ctx, task)
	if err != nil {
		return err
	}
	//新增了任务必须让该用户的“任务列表”缓存失效
	key := userTaskKey(task.UserID)
	log.Printf("[CacheStore]INVILIDATA: %s(due to CreateTask)", key)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to delete key: %s:%v", key, err)
	}
	return nil
}

// 更新
func (s *CacheStore) UpdateTask(ctx context.Context, task *models.Task) error {
	err := s.next.UpdateTask(ctx, task)
	if err != nil {
		return err
	}
	key := taskKey(task.ID)
	log.Printf("[CacheStore]INVILIDATA: %s(due to UpdateTask)", key)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to delete key: %s:%v", key, err)
	}

	//更新了任务必须让该用户的“任务列表”缓存失效
	taskKey := userTaskKey(task.UserID)
	log.Printf("[CacheStore]INVILIDATA: %s(due to UpdateTask)", key)
	if err := s.redisClient.Del(ctx, taskKey).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to delete key: %s:%v", key, err)
	}
	return nil
}

// 删除
func (s *CacheStore) DeleteTask(ctx context.Context, id int, userID int) error {
	err := s.next.DeleteTask(ctx, id, userID)
	if err != nil {
		return err
	}
	key := taskKey(id)
	log.Printf("[CacheStore]INVILIDATA: %s(due to DeleteTask)", key)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to delete key: %s:%v", key, err)
	}
	taskKey := userTaskKey(userID)
	log.Printf("[CacheStore]INVILIDATA: %s(due to DeleteTask)", key)
	if err := s.redisClient.Del(ctx, taskKey).Err(); err != nil {
		log.Printf("[CacheStore]Error: Failed to delete key: %s:%v", key, err)
	}
	return nil
}

func (s *CacheStore) CreateUser(ctx context.Context, user *models.User) error {
	return s.next.CreateUser(ctx, user)
}

func (s *CacheStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.next.GetUserByUsername(ctx, username)
}
