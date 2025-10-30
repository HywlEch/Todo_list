package middleware

import (
	"net/http"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

)

const (
	rateLimitPeriod = 1 * time.Minute //限制周期1分钟
	rateLimitMax    = 100 //周期最大请求数
)

func RateLimitMiddleware(redisClient *redis.Client)gin.HandlerFunc{
	return func(c *gin.Context){
		ctx := c.Request.Context()

		ip := c.ClientIP()
		key := "rate_limit:" + ip

		//使用redis的INCR命令为这个ip计数
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			log.Println("redis incr error:", err)
			c.Next()
			return
		}
		if count == 1 {
			redisClient.Expire(ctx, key, rateLimitPeriod)
		}
		if count > rateLimitMax{
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error":"访问过于频繁"})
		}
		c.Next()
	}
}