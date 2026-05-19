package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		//// 获取请求来源（如 http://localhost:3000）
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		//告诉浏览器：这个来源的请求我允许
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		//允许携带 Cookie / Token 等凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
