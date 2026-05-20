package controller

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

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
// 存储路径: minio://avatars/{user_id}/avatar.{ext}
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

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// ====== 删除旧头像（如果是 MinIO URL） ======
	oldUser, err := ctrl.userService.GetUserByID(userID)
	if err == nil && oldUser.Avatar != "" {
		// 如果是本地文件路径（以 /uploads/ 开头），删除本地文件
		// 如果是 MinIO URL，从 MinIO 删除
		if len(oldUser.Avatar) > 0 && oldUser.Avatar[0] == '/' {
			// 本地文件，不删除（保留兼容）
		} else if utils.MinioClient != nil {
			// 尝试从 MinIO 删除
			bucket := utils.GetBucketFromURL(oldUser.Avatar)
			if bucket != "" {
				objectName := utils.GetObjectNameFromURL(oldUser.Avatar)
				_ = utils.MinIORemoveObject(context.Background(), bucket, objectName)
			}
		}
	}

	// ====== 上传到 MinIO ======
	bucket := utils.GetMinIOConfig().BucketAvatars
	if bucket == "" {
		bucket = "avatars"
	}
	objectName := fmt.Sprintf("%d/avatar%s", userID, ext)
	contentType := "image/jpeg"
	if ext == ".png" {
		contentType = "image/png"
	} else if ext == ".webp" {
		contentType = "image/webp"
	}

	ctx := context.Background()
	if err := utils.MinIOUploadBytes(ctx, bucket, objectName, data, contentType); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 生成可访问 URL
	avatarURL := utils.MinIOGetURL(ctx, bucket, objectName)

	// 更新用户头像
	if err := ctrl.userService.UpdateAvatar(userID, avatarURL); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{"avatar": avatarURL})
}
