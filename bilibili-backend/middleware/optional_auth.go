package middleware

import (
	"strings"

	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

// OptionalAuth 可选认证中间件：解析 token 但不强制要求
// 解析成功时设置 user_id、username，失败时直接放行
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Next()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.Next()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
