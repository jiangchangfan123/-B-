package model

type Video struct {
	BaseModel
	UserID      uint64 `gorm:"not null;index" json:"user_id"`
	Title       string `gorm:"type:varchar(200);not null" json:"title"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	CoverURL    string `gorm:"type:varchar(255)" json:"cover_url"`
	VideoURL    string `gorm:"type:varchar(255)" json:"video_url"`
	Duration    int    `gorm:"default:0" json:"duration"`
	Views       int    `gorm:"default:0" json:"views"`
	Status      int    `gorm:"type:tinyint;default:1" json:"status"` // 1=已发布 2=审核中 3=已封禁
	Category    string `gorm:"type:varchar(50)" json:"category"`
}
