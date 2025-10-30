package handlers

import (
	"errors"
	"net/http"
	//"os/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/HywlEch/Todo_list/internal/models"
	"github.com/HywlEch/Todo_list/internal/store"
	"github.com/HywlEch/Todo_list/internal/config"
	"golang.org/x/crypto/bcrypt"
)

//userHandler包含用户相关得Handler
type UserHandler struct {
	Store store.Store
	JWTConfig config.JWTConfig
}
//NewUserHandler 创建一个userHandler
func NewUserHandler(s store.Store, jwtCfg config.JWTConfig)*UserHandler{
	return &UserHandler{Store: s,
		JWTConfig: jwtCfg}
}

//RegisterReqest 定义注册请求得JSON结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//Regiester处理用户注册
func (h *UserHandler)Regiester(c *gin.Context){
	var req RegisterRequest
	if err :=  c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不合理得输入"+ err.Error()})
	}
	user := &models.User{
		Username: req.Username,
		PasswordHash: req.Password,
		}
	if err := h.Store.CreateUser(c.Request.Context(), user); err != nil {
		if errors.Is(err,store.ErrUserExists){
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "用户创建成功",
		"userid": user.ID})
	
}

//LoginRequest 定义登录请求得JSON结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//LoginResponse 定义登录响应得JSON结构
type LoginResponse struct {
	Token string `json:"token"`
}

//处理用户登录
func (h *UserHandler)Login(c *gin.Context){
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不合理得输入"})
	}
	user, err := h.Store.GetUserByUsername(c.Request.Context(),req.Username) 
	if err != nil { 
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}
	//验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}
	//密码验证成功，生成JWT
	token, err := h.generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成JWT失败"})
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token})
} 

//generateJWT 生成JWT
func (h *UserHandler)generateJWT(userID int) (string, error) { 
	//定义JWT的声明
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour*time.Duration(h.JWTConfig.ExpiresInHours)).Unix(),
		"iat": time.Now().Unix, 
	}
	//创建token
	token := jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	//使用我们的密钥签名
	tokenString, err := token.SignedString([]byte(h.JWTConfig.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}