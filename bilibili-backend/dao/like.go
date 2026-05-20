package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type LikeDao struct {
	db *gorm.DB
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{db: db}
}

// Create 创建点赞记录
func (d *LikeDao) Create(like *model.VideoLike) error {
	return d.db.Create(like).Error
}

// GetByUserAndVideo 查询用户对某视频的点赞状态
func (d *LikeDao) GetByUserAndVideo(userID, videoID uint64) (*model.VideoLike, error) {
	var like model.VideoLike
	if err := d.db.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

// UpdateStatus 更新点赞状态
func (d *LikeDao) UpdateStatus(userID, videoID uint64, status int) error {
	return d.db.Model(&model.VideoLike{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Update("status", status).Error
}

// CountByVideo 统计视频点赞数
func (d *LikeDao) CountByVideo(videoID uint64) (int64, error) {
	var count int64
	err := d.db.Model(&model.VideoLike{}).
		Where("video_id = ? AND status = 1", videoID).
		Count(&count).Error
	return count, err
}
