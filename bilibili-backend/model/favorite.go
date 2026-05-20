package model

// VideoFavorite 视频收藏表
type VideoFavorite struct {
	BaseModel
	UserID  uint64 `gorm:"not null;uniqueIndex:idx_user_video_fav" json:"user_id"`
	VideoID uint64 `gorm:"not null;uniqueIndex:idx_user_video_fav" json:"video_id"`
}
