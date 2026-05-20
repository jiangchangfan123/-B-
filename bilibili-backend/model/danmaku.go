package model

import "time"

// Danmaku 弹幕模型
type Danmaku struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	VideoID   uint64    `json:"video_id" gorm:"not null;index:idx_video_time"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:varchar(100);not null"`
	TimePoint int       `json:"time_point" gorm:"default:0;index:idx_video_time"` // 视频时间点（秒）
	Color     string    `json:"color" gorm:"type:varchar(7);default:'#ffffff'"`
	Type      int8      `json:"type" gorm:"default:1"` // 1=滚动 2=顶部 3=底部
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
