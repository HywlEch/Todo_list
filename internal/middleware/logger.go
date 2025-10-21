package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc{
	return func(c *gin.Context) {
		//请求开始时间
		startTime := time.Now()

		//处理请求
		c.Next()

		//请求结束时间
		endTime := time.Now()

		//执行耗时
		latencyTime := endTime.Sub(startTime)

		//请求方式
		reqMethod := c.Request.Method

		//请求路由
		reqURI := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求IP
		clientIP := c.ClientIP()

		// 格式化日志输出
		log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqURI,
		)
	}
}