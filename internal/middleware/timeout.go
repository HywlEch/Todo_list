package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration)gin.HandlerFunc{
	return func(c *gin.Context){
		//创建一个带超时的context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()//确保在函数结束时调用cancel，释放资源

		//将新的带超时的context替换到原有的context
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		//c.Next()执行完毕之后，检查context是否超时
		if ctx.Err() == context.DeadlineExceeded{
			//如果是，设置504GatewayTimeout状态码
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
		}
	}
}