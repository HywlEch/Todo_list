package middleware
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

)

//创建一个鉴权中间件
func AuthMiddleware(jwtSecret string)gin.HandlerFunc{
	return func(c *gin.Context) {
		//从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"缺少Authorization header 字段"})
			return
		}
		//验证格式
		parts := strings.Split(authHeader, " ")
		if len(parts) !=2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header 格式必须为 'Bearer（持有人）<token>'"})
			return
		}
		tokenString := parts[1]
		//解析和验证token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		    //确保签名是我们期望的
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不正确的签名方法:%v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"验证token失败"+err.Error()})
			return
		}
		//从token中提取Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//提取用户ID并存入上下文
			userIDFloat, ok := claims["user_id"].(float64)//JWT库默认吧数字转换成浮点数
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
				return
			}
			userID := int(userIDFloat)
			//将用户ID添加到上下文
			c.Set("user_id", userID)

			//放行请求
			c.Next()
		}else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "INvalid token claims"})
			return
		}
	}
}