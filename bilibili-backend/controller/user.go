package controller

import (
	"strconv"
	"time"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService    *service.UserService
	videoService   *service.VideoService
	historyService *service.HistoryService
}

func NewUserController(userService *service.UserService, videoService *service.VideoService, historyService *service.HistoryService) *UserController {
	return &UserController{userService: userService, videoService: videoService, historyService: historyService}
}

// 获取当前用户信息 GET /api/v1/users/me
func (ctrl *UserController) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := ctrl.userService.GetUserByID(userID.(uint64))
	if err != nil {
		utils.Error(c, utils.CodeTokenInvalid)
		return
	}
	utils.OK(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"avatar":     user.Avatar,
		"sign":       user.Sign,
		"role":       user.Role,
		"coins":      user.Coins,
		"created_at": user.CreatedAt.Format(time.RFC3339),
	})
}

// 更新个人资料 PUT /api/v1/users/me
type UpdateProfileReq struct {
	Sign     string `json:"sign"`
	Nickname string `json:"nickname"`
}

func (ctrl *UserController) UpdateMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if len(req.Sign) > 200 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if len(req.Nickname) > 32 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if err := ctrl.userService.UpdateProfile(userID.(uint64), req.Sign, req.Nickname); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, nil)
}

// 修改密码 PUT /api/v1/users/me/password
type UpdatePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if req.OldPassword == "" || req.NewPassword == "" {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if len(req.NewPassword) < 6 || len(req.NewPassword) > 30 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if err := ctrl.userService.UpdatePassword(userID.(uint64), req.OldPassword, req.NewPassword); err != nil {
		if err.Error() == "旧密码错误" {
			utils.Error(c, utils.CodeUserNotFound)
			return
		}
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, nil)
}

// 获取我的视频 GET /api/v1/users/me/videos
func (ctrl *UserController) MyVideos(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	videos, total, err := ctrl.videoService.ListMyVideos(uid, page, size)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	list := make([]gin.H, 0, len(videos))
	for _, v := range videos {
		list = append(list, gin.H{
			"id":         v.ID,
			"title":      v.Title,
			"cover_url":  v.CoverURL,
			"views":      v.ViewCount,
			"status":     v.Status,
			"category":   v.Category,
			"created_at": v.CreatedAt,
		})
	}

	utils.OK(c, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// 获取播放历史 GET /api/v1/users/me/history
func (ctrl *UserController) History(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := ctrl.historyService.ListByUserID(uid, page, size)
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
