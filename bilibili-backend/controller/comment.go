package controller

import (
	"strconv"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService *service.CommentService
}

func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// Create 发表评论 POST /api/v1/videos/:id/comments
func (ctrl *CommentController) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID uint64 `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	comment, err := ctrl.commentService.CreateComment(videoID, uid, req.Content, req.ParentID)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, comment)
}

// List 获取视频评论列表 GET /api/v1/videos/:id/comments
func (ctrl *CommentController) List(c *gin.Context) {
	// 尝试获取当前登录用户ID（可选）
	var currentUserID uint64
	if uid, exists := c.Get("user_id"); exists {
		currentUserID = uid.(uint64)
	}

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	sort := c.DefaultQuery("sort", "time") // time or hot
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := ctrl.commentService.GetVideoComments(videoID, page, size, sort, currentUserID)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// Delete 删除评论 DELETE /api/v1/comments/:id
func (ctrl *CommentController) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	commentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	if err := ctrl.commentService.DeleteComment(commentID, uid); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, nil)
}

// ToggleLike 点赞/取消点赞评论 POST /api/v1/comments/:id/like
func (ctrl *CommentController) ToggleLike(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	commentID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	liked, err := ctrl.commentService.ToggleCommentLike(commentID, uid)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"liked": liked,
	})
}
