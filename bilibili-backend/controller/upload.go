package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type UploadController struct {
	userService *service.UserService
}

func NewUploadController(userService *service.UserService) *UploadController {
	return &UploadController{userService: userService}
}

// 上传头像 POST /api/v1/users/me/avatar
// 存储路径: uploads/avatars/{user_id}/avatar.{ext}
// 覆盖上传：同一用户新头像会覆盖旧文件，不保留历史版本
func (ctrl *UploadController) UploadAvatar(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := userIDVal.(uint64)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	defer file.Close()

	// 检查文件大小（5MB）
	if header.Size > 5*1024*1024 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 检查扩展名
	ext := filepath.Ext(header.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// ====== 查询并删除旧头像 ======
	oldUser, err := ctrl.userService.GetUserByID(userID)
	if err == nil && oldUser.Avatar != "" {
		// 将 URL 路径转换为文件系统路径进行删除
		oldPath := strings.TrimPrefix(oldUser.Avatar, "/")
		if oldPath != oldUser.Avatar {
			// 只删除本地存储的文件（以 /uploads/ 开头），不删除外部 URL
			_ = os.Remove(oldPath)
		}
	}

	// ====== 保存新头像（按 user_id 分组，固定文件名覆盖） ======
	uploadDir := filepath.Join("./uploads", "avatars", fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 固定文件名：avatar.{ext}，同一用户上传新头像时自动覆盖
	filename := fmt.Sprintf("avatar%s", ext)
	savePath := filepath.Join(uploadDir, filename)

	out, err := os.Create(savePath)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 生成可访问 URL
	avatarURL := fmt.Sprintf("/uploads/avatars/%d/%s", userID, filename)

	// 更新用户头像
	if err := ctrl.userService.UpdateAvatar(userID, avatarURL); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{"avatar": avatarURL})
}
