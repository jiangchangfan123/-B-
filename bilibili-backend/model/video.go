package model

// Video 视频表
type Video struct {
	BaseModel
	UserID          uint64 `gorm:"not null;index:idx_user_created" json:"user_id"`
	Title           string `gorm:"type:varchar(200);not null" json:"title"`
	Description     string `gorm:"type:text" json:"description"`
	CoverURL        string `gorm:"type:varchar(255)" json:"cover_url"`
	VideoURL        string `gorm:"type:varchar(255)" json:"video_url"`
	TranscodedURL   string `gorm:"type:varchar(255)" json:"transcoded_url"`
	Category        string `gorm:"type:varchar(50);index:idx_cat_status_created,priority:1" json:"category"`
	Status          int    `gorm:"type:tinyint;default:2;index:idx_cat_status_created,priority:2" json:"status"` // 1=已发布 2=转码中
	TranscodeStatus int    `gorm:"type:tinyint;default:0" json:"transcode_status"`                               // 0=未开始 1=转码中 2=转码完成 3=转码失败
	ViewCount       int64  `gorm:"default:0" json:"view_count"`
	LikeCount       int    `gorm:"default:0" json:"like_count"`
	CommentCount    int    `gorm:"default:0" json:"comment_count"`
	DanmakuCount    int    `gorm:"default:0" json:"danmaku_count"`
}
