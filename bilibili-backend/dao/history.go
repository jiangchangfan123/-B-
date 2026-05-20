package dao

import (
	"time"

	"bilibili-backend/model"
	"gorm.io/gorm"
)

type HistoryDao struct {
	db *gorm.DB
}

func NewHistoryDao(db *gorm.DB) *HistoryDao {
	return &HistoryDao{db: db}
}

// HistoryVideo 播放历史 JOIN 视频信息的结果
type HistoryVideo struct {
	ID        uint64    `json:"id"` // 视频 ID
	Title     string    `json:"title"`
	CoverURL  string    `json:"cover_url"`
	Views     int64     `json:"views"` // 观看数
	WatchedAt time.Time `json:"watched_at"`
}

// ListWithVideo 查询用户播放历史（JOIN 视频表）
func (d *HistoryDao) ListWithVideo(userID uint64, page, size int) ([]HistoryVideo, int64, error) {
	var list []HistoryVideo
	var total int64

	offset := (page - 1) * size

	// 统计总数
	if err := d.db.Model(&model.VideoHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// JOIN 查询
	if err := d.db.Table("video_histories").
		Select("videos.id, videos.title, videos.cover_url, videos.view_count as views, video_histories.created_at as watched_at").
		Joins("LEFT JOIN videos ON videos.id = video_histories.video_id").
		Where("video_histories.user_id = ?", userID).
		Order("video_histories.created_at DESC").
		Limit(size).Offset(offset).
		Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	if list == nil {
		list = []HistoryVideo{}
	}

	return list, total, nil
}

func (d *HistoryDao) Create(h *model.VideoHistory) error {
	return d.db.Create(h).Error
}

// First 查询单条记录
func (d *HistoryDao) First(out *model.VideoHistory, userID, videoID uint64) error {
	return d.db.Where("user_id = ? AND video_id = ?", userID, videoID).First(out).Error
}

// Save 更新记录
func (d *HistoryDao) Save(h *model.VideoHistory) error {
	return d.db.Save(h).Error
}

// CleanupOverflow 删除超出 limit 条数的最旧记录
func (d *HistoryDao) CleanupOverflow(userID uint64, limit int) error {
	var ids []uint64
	err := d.db.Model(&model.VideoHistory{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1000).
		Pluck("id", &ids).Error
	if err != nil || len(ids) <= limit {
		return nil
	}
	overflowIDs := ids[limit:]
	return d.db.Where("id IN ?", overflowIDs).Delete(&model.VideoHistory{}).Error
}

// CleanupOld 删除超过 duration 的旧记录
func (d *HistoryDao) CleanupOld(userID uint64, duration time.Duration) error {
	cutoff := time.Now().Add(-duration)
	return d.db.Where("user_id = ? AND created_at < ?", userID, cutoff).Delete(&model.VideoHistory{}).Error
}
