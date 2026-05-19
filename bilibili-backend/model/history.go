package model

type VideoHistory struct {
	BaseModel
	UserID  uint64 `gorm:"not null;index" json:"user_id"`
	VideoID uint64 `gorm:"not null;index" json:"video_id"`
}
