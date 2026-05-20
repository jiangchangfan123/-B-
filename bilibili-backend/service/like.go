package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bilibili-backend/dao"
	"bilibili-backend/model"
	"bilibili-backend/utils"
)

type LikeService struct {
	likeDao             *dao.LikeDao
	videoDao            *dao.VideoDao
	notificationService *NotificationService
}

func NewLikeService(likeDao *dao.LikeDao, videoDao *dao.VideoDao, notificationService *NotificationService) *LikeService {
	return &LikeService{likeDao: likeDao, videoDao: videoDao, notificationService: notificationService}
}

// LikeVideo 点赞
func (s *LikeService) LikeVideo(userID, videoID uint64) error {
	ctx := context.Background()

	// 1. 更新或插入 MySQL 点赞记录
	like, err := s.likeDao.GetByUserAndVideo(userID, videoID)
	if err != nil {
		// 不存在则新建
		like = &model.VideoLike{
			UserID:  userID,
			VideoID: videoID,
			Status:  1,
		}
		if err := s.likeDao.Create(like); err != nil {
			return err
		}
	} else if like.Status == 0 {
		// 已取消，重新点赞
		if err := s.likeDao.UpdateStatus(userID, videoID, 1); err != nil {
			return err
		}
	} else {
		// 已点赞，直接返回
		return nil
	}

	// 2. 更新 Redis
	if utils.RedisClient != nil {
		pipe := utils.RedisClient.Pipeline()
		pipe.HSet(ctx, fmt.Sprintf("video:likes:user:%d", userID), strconv.FormatUint(videoID, 10), 1)
		pipe.Incr(ctx, fmt.Sprintf("video:likes:count:%d", videoID))
		pipe.SAdd(ctx, "video:likes:pending", videoID)
		pipe.Expire(ctx, fmt.Sprintf("video:likes:user:%d", userID), 24*time.Hour)
		pipe.Expire(ctx, fmt.Sprintf("video:likes:count:%d", videoID), 24*time.Hour)
		_, _ = pipe.Exec(ctx)
	}

	// 3. 发送通知（异步，不阻塞）
	if s.notificationService != nil {
		video, _ := s.videoDao.GetByID(videoID)
		if video != nil && video.UserID != userID {
			go s.notificationService.CreateVideoLikedNotification(videoID, userID, video.UserID)
		}
	}

	return nil
}

// UnlikeVideo 取消点赞
func (s *LikeService) UnlikeVideo(userID, videoID uint64) error {
	ctx := context.Background()

	like, err := s.likeDao.GetByUserAndVideo(userID, videoID)
	if err != nil || like.Status == 0 {
		return nil // 未点赞或不存在
	}

	if err := s.likeDao.UpdateStatus(userID, videoID, 0); err != nil {
		return err
	}

	if utils.RedisClient != nil {
		pipe := utils.RedisClient.Pipeline()
		pipe.HSet(ctx, fmt.Sprintf("video:likes:user:%d", userID), strconv.FormatUint(videoID, 10), 0)
		pipe.Decr(ctx, fmt.Sprintf("video:likes:count:%d", videoID))
		pipe.SAdd(ctx, "video:likes:pending", videoID)
		_, _ = pipe.Exec(ctx)
	}

	return nil
}

// ToggleLike 点赞/取消 toggle
func (s *LikeService) ToggleLike(userID, videoID uint64) (liked bool, err error) {
	like, err := s.likeDao.GetByUserAndVideo(userID, videoID)
	if err != nil || like.Status == 0 {
		// 未点赞，执行点赞
		if err := s.LikeVideo(userID, videoID); err != nil {
			return false, err
		}
		return true, nil
	}
	// 已点赞，执行取消
	if err := s.UnlikeVideo(userID, videoID); err != nil {
		return true, err
	}
	return false, nil
}

// IsLiked 查询用户是否点赞
func (s *LikeService) IsLiked(userID, videoID uint64) (bool, error) {
	ctx := context.Background()

	// 先查 Redis
	if utils.RedisClient != nil {
		val, err := utils.RedisClient.HGet(ctx, fmt.Sprintf("video:likes:user:%d", userID), strconv.FormatUint(videoID, 10)).Result()
		if err == nil {
			return val == "1", nil
		}
	}

	// 再查 MySQL
	like, err := s.likeDao.GetByUserAndVideo(userID, videoID)
	if err != nil {
		return false, nil
	}
	return like.Status == 1, nil
}

// GetLikeCount 获取点赞数（优先 Redis）
func (s *LikeService) GetLikeCount(videoID uint64) (int64, error) {
	ctx := context.Background()

	if utils.RedisClient != nil {
		val, err := utils.RedisClient.Get(ctx, fmt.Sprintf("video:likes:count:%d", videoID)).Result()
		if err == nil {
			count, _ := strconv.ParseInt(val, 10, 64)
			return count, nil
		}
	}

	return s.likeDao.CountByVideo(videoID)
}

// SyncLikesToDB 将 Redis 中的点赞数同步回 MySQL（定时任务调用）
func (s *LikeService) SyncLikesToDB() error {
	if utils.RedisClient == nil {
		return nil
	}
	ctx := context.Background()

	videoIDs, err := utils.RedisClient.SMembers(ctx, "video:likes:pending").Result()
	if err != nil || len(videoIDs) == 0 {
		return nil
	}

	for _, vidStr := range videoIDs {
		videoID, _ := strconv.ParseUint(vidStr, 10, 64)
		count, _ := s.likeDao.CountByVideo(videoID)
		_ = s.videoDao.UpdateLikeCount(videoID, int(count))
	}

	_ = utils.RedisClient.Del(ctx, "video:likes:pending")
	return nil
}
