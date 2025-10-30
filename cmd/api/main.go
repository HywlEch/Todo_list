// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HywlEch/Todo_list/internal/config"
	"github.com/HywlEch/Todo_list/internal/handlers"
	"github.com/HywlEch/Todo_list/internal/middleware"
	"github.com/HywlEch/Todo_list/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func initRedisClient(cfg config.RedisConfig)*redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
		DB: cfg.DB,
	})
	//检查连接
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Redis连接失败：%s",err)
	}
	log.Println("Redis连接成功")
	return rdb
}
func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}
	//初始化Store
	dbStore, err := store.NewPostgresStore(cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败：%s",err)
	}
	defer dbStore.DB.Close()

	//初始化Redis客户端
	redisClient := initRedisClient(cfg.Redis)
	defer redisClient.Close()

	//初始化Handler
	taskHandler := handlers.NewTaskHandler(dbStore)
	//初始化UserHandler，传入JWT配置
	userHandler := handlers.NewUserHandler(dbStore, cfg.JWT)

	//设置路由

	//router := gin.Default()
	router := gin.New() //使用gin.New()创建一个干净的引擎
	//全局应用中间件
	router.Use(middleware.Logger()) //应用日志中间件
	router.Use(gin.Recovery()) //使用gin默认的Recovery中间件,防止panic
	router.Use(middleware.RateLimitMiddleware(redisClient))
	router.Use(middleware.TimeoutMiddleware(10 * time.Second))//应用5秒钟超时中间件

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/regist", userHandler.Regiester)
		authRouter.POST("/login", userHandler.Login)
	}

	taskRouter := router.Group("/tasks")
	{
		taskRouter.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		taskRouter.POST("", taskHandler.CreateTask)
		taskRouter.GET("", taskHandler.GetTasks)
		taskRouter.GET("/:id", taskHandler.GetTaskByID)
		taskRouter.PUT("/:id", taskHandler.UpdateTask)
		taskRouter.DELETE("/:id", taskHandler.DeleteTask)
	}

	// serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	// log.Printf("Server is running on port %s...", cfg.Server.Port)
	// if err := router.Run(serverAddr); err != nil {
	// 	log.Fatalf("Could not start server: %s\n", err)
	// }
	//创建一个http.Server实力，提供更多的控制权
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	//在一个GoRoutine中启动服务器，这样他就不会阻塞主线程
	go func(){
		log.Printf("服务正在运行...:%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//创建一个channel来接受操作系统的中断信号
	quit := make(chan os.Signal, 1)
	//signal.Notify函数会将指定得信号发送到quit channel
	//在这里监听SIGINT(CTRL+C)和SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//阻塞主进程,直到quitchannel 收到一个信号
	<-quit
	log.Println("关闭服务中...")

	//创造一个有超时的context， 用于通知我们有5秒的时间来处理请求
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()

	//调用Shutdown函数来关闭服务器
	//停止接受新的请求并等待现有请求
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}