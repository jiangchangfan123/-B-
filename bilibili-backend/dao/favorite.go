package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type FavoriteDao struct {
	db *gorm.DB
}

func NewFavoriteDao(db *gorm.DB) *FavoriteDao {
	return &FavoriteDao{db: db}
}

// Create 收藏
func (d *FavoriteDao) Create(f *model.VideoFavorite) error {
	return d.db.Create(f).Error
}

// Delete 取消收藏
func (d *FavoriteDao) Delete(userID, videoID uint64) error {
	return d.db.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&model.VideoFavorite{}).Error
}

// GetByUserAndVideo 查询是否收藏
func (d *FavoriteDao) GetByUserAndVideo(userID, videoID uint64) (*model.VideoFavorite, error) {
	var f model.VideoFavorite
	if err := d.db.Where("user_id = ? AND video_id = ?", userID, videoID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

// FavoriteVideo 收藏视频 JOIN 查询结果
type FavoriteVideo struct {
	VideoID   uint64 `json:"video_id"`
	Title     string `json:"title"`
	CoverURL  string `json:"cover_url"`
	Category  string `json:"category"`
	ViewCount int64  `json:"view_count"`
	LikeCount int    `json:"like_count"`
}

// ListWithVideo 查询用户收藏列表（JOIN 视频表）
func (d *FavoriteDao) ListWithVideo(userID uint64, page, size int) ([]FavoriteVideo, int64, error) {
	var list []FavoriteVideo
	var total int64

	offset := (page - 1) * size

	if err := d.db.Model(&model.VideoFavorite{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.Table("video_favorites").
		Select("videos.id as video_id, videos.title, videos.cover_url, videos.category, videos.view_count, videos.like_count").
		Joins("LEFT JOIN videos ON videos.id = video_favorites.video_id").
		Where("video_favorites.user_id = ?", userID).
		Order("video_favorites.created_at DESC").
		Limit(size).Offset(offset).
		Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
