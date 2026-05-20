package controller

import (
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type VideoController struct {
	videoService   *service.VideoService
	userService    *service.UserService
	historyService *service.HistoryService
}

func NewVideoController(videoService *service.VideoService, userService *service.UserService, historyService *service.HistoryService) *VideoController {
	return &VideoController{videoService: videoService, userService: userService, historyService: historyService}
}

// Upload 视频上传 POST /api/v1/videos
func (ctrl *VideoController) Upload(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	title := c.PostForm("title")
	if title == "" || len(title) > 200 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	description := c.PostForm("description")
	category := c.PostForm("category")
	if category == "" {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 视频文件
	videoFile, videoHeader, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	defer videoFile.Close()

	// 校验格式
	videoExt := strings.ToLower(filepath.Ext(videoHeader.Filename))
	allowedExts := map[string]bool{".mp4": true, ".mov": true, ".webm": true, ".mkv": true}
	if !allowedExts[videoExt] {
		c.JSON(400, gin.H{"code": 40003, "message": "视频格式不支持", "data": nil})
		return
	}
	// 大小限制 2GB
	if videoHeader.Size > 2*1024*1024*1024 {
		c.JSON(400, gin.H{"code": 40004, "message": "文件过大，最大支持 2GB", "data": nil})
		return
	}

	// 封面文件（可选）
	var coverFile io.Reader
	var coverExt string
	coverFileObj, coverHeader, err := c.Request.FormFile("cover")
	if err == nil && coverFileObj != nil {
		defer coverFileObj.Close()
		coverExt = strings.ToLower(filepath.Ext(coverHeader.Filename))
		if coverExt != ".jpg" && coverExt != ".jpeg" && coverExt != ".png" && coverExt != ".webp" {
			coverFile = nil
			coverExt = ""
		} else {
			coverFile = coverFileObj
		}
	}

	video, err := ctrl.videoService.UploadVideo(uid, title, description, category, videoFile, coverFile, videoExt, coverExt, videoHeader.Size)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"id":               video.ID,
		"title":            video.Title,
		"cover_url":        video.CoverURL,
		"status":           video.Status,
		"transcode_status": video.TranscodeStatus,
		"created_at":       video.CreatedAt,
	})
}

// TranscodeStatus 转码状态 GET /api/v1/videos/:id/transcode
func (ctrl *VideoController) TranscodeStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	video, err := ctrl.videoService.GetVideoByID(id)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, gin.H{
		"transcode_status": video.TranscodeStatus,
		"transcoded_url":   video.TranscodedURL,
	})
}

// List 视频列表 GET /api/v1/videos
func (ctrl *VideoController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	category := c.Query("category")
	sort := c.DefaultQuery("sort", "new")

	videos, total, err := ctrl.videoService.ListPublishedVideos(category, sort, page, size)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 查询 UP 主信息
	userIDs := make(map[uint64]bool)
	for _, v := range videos {
		userIDs[v.UserID] = true
	}
	userMap := make(map[uint64]*gin.H)
	for uid := range userIDs {
		user, _ := ctrl.userService.GetUserByID(uid)
		if user != nil {
			userMap[uid] = &gin.H{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.Nickname,
				"avatar":   user.Avatar,
			}
		}
	}

	list := make([]gin.H, 0, len(videos))
	for _, v := range videos {
		item := gin.H{
			"id":         v.ID,
			"title":      v.Title,
			"cover_url":  v.CoverURL,
			"duration":   v.Duration,
			"category":   v.Category,
			"view_count": v.ViewCount,
			"like_count": v.LikeCount,
			"created_at": v.CreatedAt,
			"user":       userMap[v.UserID],
		}
		list = append(list, item)
	}

	utils.OK(c, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// Detail 视频详情 GET /api/v1/videos/:id
func (ctrl *VideoController) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	video, err := ctrl.videoService.GetVideoDetail(id)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 播放量 +1
	_ = ctrl.videoService.IncrementView(id)

	// 记录播放历史（如果已登录）
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint64)
		_ = ctrl.historyService.RecordHistory(uid, id)
	}

	// UP 主信息
	user, _ := ctrl.userService.GetUserByID(video.UserID)
	userInfo := gin.H{}
	if user != nil {
		userInfo = gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		}
	}

	utils.OK(c, gin.H{
		"id":               video.ID,
		"title":            video.Title,
		"description":      video.Description,
		"cover_url":        video.CoverURL,
		"video_url":        video.VideoURL,
		"transcoded_url":   video.TranscodedURL,
		"duration":         video.Duration,
		"category":         video.Category,
		"status":           video.Status,
		"transcode_status": video.TranscodeStatus,
		"view_count":       video.ViewCount,
		"like_count":       video.LikeCount,
		"comment_count":    video.CommentCount,
		"danmaku_count":    video.DanmakuCount,
		"created_at":       video.CreatedAt,
		"user_info":        userInfo,
	})
}

// Delete 删除视频 DELETE /api/v1/videos/:id
func (ctrl *VideoController) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	if err := ctrl.videoService.DeleteVideo(uid, id); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, nil)
}

// Update 更新视频 PUT /api/v1/videos/:id
func (ctrl *VideoController) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	category := c.PostForm("category")

	if title == "" || len(title) > 200 {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 可选的视频文件
	var videoFile io.Reader
	var videoExt string
	file, header, err := c.Request.FormFile("video")
	if err == nil && header != nil {
		videoFile = file
		videoExt = filepath.Ext(header.Filename)
		defer file.Close()
	}

	// 可选的封面
	var coverFile io.Reader
	var coverExt string
	cover, coverHeader, err := c.Request.FormFile("cover")
	if err == nil && coverHeader != nil {
		coverFile = cover
		coverExt = filepath.Ext(coverHeader.Filename)
		defer cover.Close()
	}

	if err := ctrl.videoService.UpdateVideo(uid, videoID, title, description, category, videoFile, coverFile, videoExt, coverExt); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, nil)
}
