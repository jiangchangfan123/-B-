package model

// Comment 评论表
type Comment struct {
	BaseModel
	VideoID   uint64 `gorm:"not null;index:idx_video_parent" json:"video_id"`
	UserID    uint64 `gorm:"not null;index" json:"user_id"`
	Content   string `gorm:"type:text;not null" json:"content"`
	ParentID  uint64 `gorm:"default:0;index:idx_video_parent" json:"parent_id"` // 回复的评论ID，0=一级评论
	RootID    uint64 `gorm:"default:0;index" json:"root_id"`                    // 一级评论ID，方便聚合回复
	LikeCount int    `gorm:"default:0" json:"like_count"`
	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
