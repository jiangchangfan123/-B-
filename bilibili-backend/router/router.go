package router

import (
	"bilibili-backend/controller"
	"bilibili-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authCtrl *controller.AuthController, userCtrl *controller.UserController, uploadCtrl *controller.UploadController, videoCtrl *controller.VideoController, likeCtrl *controller.LikeController, favoriteCtrl *controller.FavoriteController, commentCtrl *controller.CommentController, notificationCtrl *controller.NotificationController) {
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
			users.GET("/me/notifications", notificationCtrl.List)
			users.GET("/me/notifications/unread-count", notificationCtrl.UnreadCount)
			users.PUT("/me/notifications/:id/read", notificationCtrl.MarkAsRead)
			users.PUT("/me/notifications/read-all", notificationCtrl.MarkAllAsRead)
			users.DELETE("/me/notifications/:id", notificationCtrl.Delete)
		}

		// 视频组（部分需 JWT）
		videos := apiV1.Group("/videos")
		{
			videos.GET("", videoCtrl.List)
			videos.GET("/:id/transcode", videoCtrl.TranscodeStatus)
		}
		// 视频详情（可选认证，已登录用户记录播放历史）
		videosDetail := apiV1.Group("/videos")
		videosDetail.Use(middleware.OptionalAuth())
		{
			videosDetail.GET("/:id", videoCtrl.Detail)
			videosDetail.GET("/:id/comments", commentCtrl.List)
		}
		// 需登录的视频操作
		authVideos := apiV1.Group("/videos")
		authVideos.Use(middleware.Auth())
		{
			authVideos.POST("", videoCtrl.Upload)
			authVideos.DELETE("/:id", videoCtrl.Delete)
			authVideos.PUT("/:id", videoCtrl.Update)
			authVideos.POST("/:id/like", likeCtrl.ToggleLike)
			authVideos.GET("/:id/like/status", likeCtrl.LikeStatus)
			authVideos.POST("/:id/favorite", favoriteCtrl.ToggleFavorite)
			authVideos.GET("/:id/favorite/status", favoriteCtrl.FavoriteStatus)
			authVideos.POST("/:id/comments", commentCtrl.Create)
		}

		// 评论相关
		authComments := apiV1.Group("/comments")
		authComments.Use(middleware.Auth())
		{
			authComments.DELETE("/:id", commentCtrl.Delete)
			authComments.POST("/:id/like", commentCtrl.ToggleLike)
		}

		// 收藏列表
		users.GET("/me/favorites", favoriteCtrl.ListFavorites)
	}
}
