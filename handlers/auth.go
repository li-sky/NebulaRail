package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/li-sky/NebulaRail/config" // 使用 go.mod 中的模块路径
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求结构体
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// LoginHandler 处理登录请求
func LoginHandler(appConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginReq LoginRequest

		if err := c.BindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
			return
		}

		// 验证密码
		if err := bcrypt.CompareHashAndPassword([]byte(appConfig.PasswordHash), []byte(loginReq.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
			return
		}

		// 生成JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"authorized": true,
			"exp":        time.Now().Add(appConfig.JWTExpiry).Unix(),
		})

		tokenString, err := token.SignedString([]byte(appConfig.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":      tokenString,
			"expires_in": appConfig.JWTExpiry.Seconds(),
		})
	}
}

// ChangePasswordHandler 处理修改密码请求
func ChangePasswordHandler(appConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePasswordRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
			return
		}

		// 验证当前密码
		if err := bcrypt.CompareHashAndPassword([]byte(appConfig.PasswordHash), []byte(req.CurrentPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "当前密码错误"})
			return
		}

		// 生成新密码hash
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理错误"})
			return
		}

		// 将新密码哈希保存到数据库
		if err := config.SetConfigValue(appConfig.DB, "password_hash", string(newPasswordHash)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存新密码"})
			return
		}

		// 更新内存中的配置 (重要!)
		appConfig.PasswordHash = string(newPasswordHash)

		c.JSON(http.StatusOK, gin.H{"message": "密码已更新"})
	}
}

// ProtectedHandler 示例受保护路由
func ProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "你已通过认证"})
	}
}
