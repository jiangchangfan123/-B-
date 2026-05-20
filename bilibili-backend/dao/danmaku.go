package dao

import (
	"bilibili-backend/model"

	"gorm.io/gorm"
)

// DanmakuDAO 弹幕数据访问层
type DanmakuDAO struct{}

// Create 创建弹幕记录
func (d *DanmakuDAO) Create(db *gorm.DB, danmaku *model.Danmaku) error {
	return db.Create(danmaku).Error
}

// GetByVideoIDTimeRange 获取某视频某时间段的弹幕
func (d *DanmakuDAO) GetByVideoIDTimeRange(db *gorm.DB, videoID uint64, start, end int) ([]model.Danmaku, error) {
	var danmakuList []model.Danmaku
	err := db.Where("video_id = ? AND time_point >= ? AND time_point <= ?", videoID, start, end).
		Order("time_point ASC, created_at ASC").
		Find(&danmakuList).Error
	return danmakuList, err
}
