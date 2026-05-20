package controller

import (
	"strconv"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type LikeController struct {
	likeService *service.LikeService
}

func NewLikeController(likeService *service.LikeService) *LikeController {
	return &LikeController{likeService: likeService}
}

// ToggleLike POST /api/v1/videos/:id/like
func (ctrl *LikeController) ToggleLike(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	liked, err := ctrl.likeService.ToggleLike(uid, videoID)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 获取最新点赞数
	count, _ := ctrl.likeService.GetLikeCount(videoID)

	utils.OK(c, gin.H{
		"liked": liked,
		"count": count,
	})
}

// LikeStatus GET /api/v1/videos/:id/like/status
func (ctrl *LikeController) LikeStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	liked, _ := ctrl.likeService.IsLiked(uid, videoID)
	count, _ := ctrl.likeService.GetLikeCount(videoID)

	utils.OK(c, gin.H{
		"liked": liked,
		"count": count,
	})
}
