package router

import (
	"bilibili-backend/controller"
	"bilibili-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authCtrl *controller.AuthController, userCtrl *controller.UserController, uploadCtrl *controller.UploadController, videoCtrl *controller.VideoController) {
	// 全局中间件
	r.Use(middleware.CORS())

	// 静态文件服务（头像上传目录）
	r.Static("/uploads", "./uploads")

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		// 认证组（无需 JWT）
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
		}

		// 用户组（需 JWT）
		users := apiV1.Group("/users")
		users.Use(middleware.Auth())
		{
			users.GET("/me", userCtrl.Me)
			users.PUT("/me", userCtrl.UpdateMe)
			users.PUT("/me/password", userCtrl.UpdatePassword)
			users.POST("/me/avatar", uploadCtrl.UploadAvatar)
			users.GET("/me/videos", userCtrl.MyVideos)
			users.GET("/me/history", userCtrl.History)
		}

		// 视频组（部分需 JWT）
		videos := apiV1.Group("/videos")
		{
			videos.GET("", videoCtrl.List)
			videos.GET("/:id", videoCtrl.Detail)
			videos.GET("/:id/transcode", videoCtrl.TranscodeStatus)
		}
		// 需登录的视频操作
		authVideos := apiV1.Group("/videos")
		authVideos.Use(middleware.Auth())
		{
			authVideos.POST("", videoCtrl.Upload)
			authVideos.DELETE("/:id", videoCtrl.Delete)
		}
	}
}
