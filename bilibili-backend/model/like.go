package model

// VideoLike 视频点赞表
type VideoLike struct {
	BaseModel
	UserID  uint64 `gorm:"not null;uniqueIndex:idx_user_video_like" json:"user_id"`
	VideoID uint64 `gorm:"not null;uniqueIndex:idx_user_video_like" json:"video_id"`
	Status  int    `gorm:"default:1" json:"status"` // 1=点赞 0=取消
}
