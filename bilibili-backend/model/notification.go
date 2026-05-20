package model

// NotificationType 通知类型
const (
	NotificationTypeCommentReply = 1 // 评论被回复
	NotificationTypeCommentLiked = 2 // 评论被点赞
	NotificationTypeVideoLiked   = 3 // 视频被点赞
)

// Notification 消息通知表
type Notification struct {
	BaseModel
	UserID        uint64 `gorm:"not null;index:idx_user_read" json:"user_id"` // 接收者
	Type          int    `gorm:"not null" json:"type"`                        // 1=评论回复 2=评论被点赞 3=视频被点赞
	Title         string `gorm:"size:100;not null" json:"title"`              // 通知标题
	Content       string `gorm:"size:255" json:"content"`                     // 通知内容
	RelatedID     uint64 `gorm:"index" json:"related_id"`                     // 关联ID（视频ID或评论ID）
	IsRead        bool   `gorm:"default:false;index:idx_user_read" json:"is_read"`
	TriggerUserID uint64 `gorm:"" json:"trigger_user_id"` // 触发用户ID
	TriggerUser   User   `gorm:"foreignKey:TriggerUserID" json:"trigger_user,omitempty"`
}
