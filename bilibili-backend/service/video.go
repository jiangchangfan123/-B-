package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bilibili-backend/config"
	"bilibili-backend/dao"
	"bilibili-backend/model"
	"bilibili-backend/utils"
	"github.com/google/uuid"
)

type VideoService struct {
	videoDao *dao.VideoDao
	userDao  *dao.UserDao
}

func NewVideoService(videoDao *dao.VideoDao, userDao *dao.UserDao) *VideoService {
	return &VideoService{videoDao: videoDao, userDao: userDao}
}

// UploadVideo 视频上传（含封面处理、MinIO 上传、发送转码任务）
func (s *VideoService) UploadVideo(userID uint64, title, description, category string, videoFile, coverFile io.Reader, videoExt, coverExt string, videoSize int64) (*model.Video, error) {
	videoIDStr := uuid.New().String()

	// === 1. 处理视频上传 ===
	// 先保存到本地临时文件
	tmpDir := filepath.Join(os.TempDir(), "uploads", videoIDStr)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir) // 上传完成后清理

	videoTmpPath := filepath.Join(tmpDir, "original"+videoExt)
	if err := saveReaderToFile(videoFile, videoTmpPath); err != nil {
		return nil, fmt.Errorf("save video failed: %w", err)
	}

	// 上传原片到 MinIO
	minioVideoKey := fmt.Sprintf("videos/raw/%s/original%s", videoIDStr, videoExt)
	if err := utils.MinIOUploadFile(context.Background(), config.C.MinIO.BucketVideos, minioVideoKey, videoTmpPath, "video/mp4"); err != nil {
		return nil, fmt.Errorf("upload video to minio failed: %w", err)
	}
	videoURL := utils.MinIOGetURL(context.Background(), config.C.MinIO.BucketVideos, minioVideoKey)

	// === 2. 处理封面 ===
	var coverURL string
	if coverFile != nil && coverExt != "" {
		// 用户上传封面
		coverTmpPath := filepath.Join(tmpDir, "cover"+coverExt)
		if err := saveReaderToFile(coverFile, coverTmpPath); err != nil {
			return nil, fmt.Errorf("save cover failed: %w", err)
		}
		minioCoverKey := fmt.Sprintf("covers/%s%s", videoIDStr, coverExt)
		if err := utils.MinIOUploadFile(context.Background(), config.C.MinIO.BucketCovers, minioCoverKey, coverTmpPath, "image/jpeg"); err != nil {
			return nil, fmt.Errorf("upload cover to minio failed: %w", err)
		}
		coverURL = utils.MinIOGetURL(context.Background(), config.C.MinIO.BucketCovers, minioCoverKey)
	} else {
		// FFmpeg 截取第 5 秒作为封面
		coverTmpPath := filepath.Join(tmpDir, "cover.jpg")
		if _, err := exec.LookPath("ffmpeg"); err == nil {
			cmd := exec.Command("ffmpeg", "-ss", "00:00:05", "-i", videoTmpPath, "-vframes", "1", "-q:v", "2", coverTmpPath)
			if err := cmd.Run(); err == nil {
				minioCoverKey := fmt.Sprintf("covers/%s.jpg", videoIDStr)
				if err := utils.MinIOUploadFile(context.Background(), config.C.MinIO.BucketCovers, minioCoverKey, coverTmpPath, "image/jpeg"); err == nil {
					coverURL = utils.MinIOGetURL(context.Background(), config.C.MinIO.BucketCovers, minioCoverKey)
				}
			}
		}
	}

	// === 3. 写入数据库 ===
	video := &model.Video{
		UserID:          userID,
		Title:           title,
		Description:     description,
		CoverURL:        coverURL,
		VideoURL:        videoURL,
		Category:        category,
		Status:          2, // 转码中
		TranscodeStatus: 1, // 转码中
	}
	if err := s.videoDao.Create(video); err != nil {
		return nil, fmt.Errorf("create video record failed: %w", err)
	}

	// === 4. 发送转码任务到 RabbitMQ ===
	transcodeMsg := utils.TranscodeMessage{
		VideoID:   video.ID,
		InputURL:  videoURL,
		Bucket:    config.C.MinIO.BucketVideos,
		ObjectKey: fmt.Sprintf("videos/480p/%s_480p.mp4", videoIDStr),
	}
	if err := utils.PublishTranscodeTask(transcodeMsg); err != nil {
		// 发送失败不阻断上传，只记录日志
		fmt.Printf("[WARN] publish transcode task failed: %v\n", err)
	}

	return video, nil
}

// GetVideoDetail 获取视频详情
func (s *VideoService) GetVideoDetail(id uint64) (*model.Video, error) {
	return s.videoDao.GetByID(id)
}

// IncrementView 播放量 +1
func (s *VideoService) IncrementView(id uint64) error {
	return s.videoDao.IncrementViewCount(id)
}

// DeleteVideo 删除视频（检验所有权 + 删除 MinIO 文件）
func (s *VideoService) DeleteVideo(userID, videoID uint64) error {
	video, err := s.videoDao.GetByID(videoID)
	if err != nil {
		return err
	}
	if video.UserID != userID {
		return fmt.Errorf("not owner")
	}

	// 删除 MinIO 上的文件
	ctx := context.Background()
	if video.VideoURL != "" {
		key := extractObjectKey(video.VideoURL)
		if key != "" {
			_ = utils.MinIORemoveObject(ctx, config.C.MinIO.BucketVideos, key)
		}
	}
	if video.TranscodedURL != "" {
		key := extractObjectKey(video.TranscodedURL)
		if key != "" {
			_ = utils.MinIORemoveObject(ctx, config.C.MinIO.BucketVideos, key)
		}
	}
	if video.CoverURL != "" {
		key := extractObjectKey(video.CoverURL)
		if key != "" {
			_ = utils.MinIORemoveObject(ctx, config.C.MinIO.BucketCovers, key)
		}
	}

	return s.videoDao.Delete(videoID)
}

// ListPublishedVideos 已发布视频列表
func (s *VideoService) ListPublishedVideos(category, sort string, page, size int) ([]model.Video, int64, error) {
	return s.videoDao.ListPublished(category, sort, page, size)
}

// ListMyVideos 我的视频列表
func (s *VideoService) ListMyVideos(userID uint64, page, size int) ([]model.Video, int64, error) {
	return s.videoDao.ListByUserID(userID, page, size)
}

// GetVideoByID 根据 ID 获取视频
func (s *VideoService) GetVideoByID(id uint64) (*model.Video, error) {
	return s.videoDao.GetByID(id)
}

// UpdateTranscodeResult 更新转码结果
func (s *VideoService) UpdateTranscodeResult(videoID uint64, transcodedURL string, success bool) error {
	if success {
		return s.videoDao.UpdateTranscodeStatus(videoID, 2, transcodedURL, 1)
	}
	return s.videoDao.UpdateTranscodeStatus(videoID, 3, "", 2)
}

// 辅助函数
func saveReaderToFile(r io.Reader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

// 从 MinIO URL 中提取 object key
func extractObjectKey(url string) string {
	// 格式: http://endpoint/bucket/objectKey
	parts := strings.SplitN(url, "/", 5)
	if len(parts) < 5 {
		return ""
	}
	return parts[4]
}
