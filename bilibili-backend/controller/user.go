package controller

import (
	"strconv"
	"time"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
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

	// Mock 视频数据（视频功能尚未实现）
	mockVideos := []gin.H{
		{"id": 1, "title": "深空漫步者 // 第一章", "cover_url": "", "views": 3420, "status": 1, "category": "科幻", "created_at": time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
		{"id": 2, "title": "虚空协议 // 实验记录", "cover_url": "", "views": 1280, "status": 2, "category": "技术", "created_at": time.Now().Add(-72 * time.Hour).Format(time.RFC3339)},
		{"id": 3, "title": "星际通信 // 第三节", "cover_url": "", "views": 8900, "status": 1, "category": "游戏", "created_at": time.Now().Add(-168 * time.Hour).Format(time.RFC3339)},
	}

	total := int64(len(mockVideos))
	utils.OK(c, gin.H{
		"list":  mockVideos,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// 获取播放历史 GET /api/v1/users/me/history
func (ctrl *UserController) History(c *gin.Context) {
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

	// Mock 播放历史（历史功能尚未实现）
	mockHistory := []gin.H{
		{"id": 1, "title": "深空漫步者 // 第一章", "cover_url": "", "views": 3420, "watched_at": time.Now().Add(-2 * time.Hour).Format(time.RFC3339)},
		{"id": 2, "title": "虚空协议 // 实验记录", "cover_url": "", "views": 1280, "watched_at": time.Now().Add(-5 * time.Hour).Format(time.RFC3339)},
		{"id": 4, "title": "星际旅行 // 红巨星探索", "cover_url": "", "views": 5600, "watched_at": time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
	}

	total := int64(len(mockHistory))
	utils.OK(c, gin.H{
		"list":  mockHistory,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
