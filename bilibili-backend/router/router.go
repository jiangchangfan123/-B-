package router

import (
	"bilibili-backend/controller"
	"bilibili-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authCtrl *controller.AuthController) {
	// 全局中间件
	r.Use(middleware.CORS())

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
			auth.GET("/me", middleware.Auth(), authCtrl.Me)
		}
	}
}
