package middleware

import (
	"errors"
	"net/http"
	"log"

	"github.com/HywlEch/Todo_list/internal/store"
	"github.com/HywlEch/Todo_list/internal/apperrors"
	"github.com/gin-gonic/gin"
)

//ErrorMiddleware 全局错误中间件
func ErrorMiddleware() gin.HandlerFunc { 
	return func(c *gin.Context){
		//首先执行链路中的下一个handler
		c.Next()

		//c.Error()是一个[]*gin.Error
		if len(c.Errors) == 0 { 
			return
		}

		err := c.Errors[0].Err

		//定义默认的错误响应
		httpCode := http.StatusInternalServerError
		jsonResponse := gin.H{"errors": "Internal Server Error"}

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) { 
			httpCode = appErr.Code
			jsonResponse = gin.H{"errors": appErr.Message}
		}else if errors.Is(err, store.ErrNotFound) { 
			httpCode = http.StatusNotFound
			jsonResponse = gin.H{"errors": "Not Found"}
		}else if errors.Is(err, store.ErrUserExists) { 
			httpCode = http.StatusConflict
			jsonResponse = gin.H{"errors": "User Already Exists"}
		}

		//记录日志 500错误需要记录完整得错误信息，而4XX错误只需要info级别
		if httpCode >=500 {
			log.Printf("Internal Server Error: %v\nFull error chain: %+v", err, err)
		}else {
			log.Printf("Client Error(%d): %v", httpCode, err)
		}
		//返回json响应
	if !c.Writer.Written() {
		c.AbortWithStatusJSON(httpCode, jsonResponse)
	}
	}
}