package service

import (
	"fmt"

	"bilibili-backend/dao"
	"bilibili-backend/model"
)

type NotificationService struct {
	notificationDao *dao.NotificationDao
	userDao         *dao.UserDao
	videoDao        *dao.VideoDao
	commentDao      *dao.CommentDao
}

func NewNotificationService(notificationDao *dao.NotificationDao, userDao *dao.UserDao, videoDao *dao.VideoDao, commentDao *dao.CommentDao) *NotificationService {
	return &NotificationService{
		notificationDao: notificationDao,
		userDao:         userDao,
		videoDao:        videoDao,
		commentDao:      commentDao,
	}
}

// CreateCommentReplyNotification 评论被回复
func (s *NotificationService) CreateCommentReplyNotification(videoID, commentID, replyUserID, toUserID uint64, content string) {
	if replyUserID == toUserID {
		return // 自己回复自己不发通知
	}
	triggerUser, _ := s.userDao.GetByID(replyUserID)
	nickname := "某人"
	if triggerUser != nil {
		nickname = triggerUser.Nickname
		if nickname == "" {
			nickname = triggerUser.Username
		}
	}
	_ = s.notificationDao.Create(&model.Notification{
		UserID:        toUserID,
		Type:          model.NotificationTypeCommentReply,
		Title:         fmt.Sprintf("%s 回复了你的评论", nickname),
		Content:       content,
		RelatedID:     commentID,
		TriggerUserID: replyUserID,
	})
}

// CreateCommentLikedNotification 评论被点赞
func (s *NotificationService) CreateCommentLikedNotification(commentID, likeUserID, toUserID uint64) {
	if likeUserID == toUserID {
		return
	}
	triggerUser, _ := s.userDao.GetByID(likeUserID)
	nickname := "某人"
	if triggerUser != nil {
		nickname = triggerUser.Nickname
		if nickname == "" {
			nickname = triggerUser.Username
		}
	}
	comment, _ := s.commentDao.GetByID(commentID)
	content := ""
	if comment != nil {
		content = comment.Content
		if len(content) > 50 {
			content = content[:50] + "..."
		}
	}
	_ = s.notificationDao.Create(&model.Notification{
		UserID:        toUserID,
		Type:          model.NotificationTypeCommentLiked,
		Title:         fmt.Sprintf("%s 点赞了你的评论", nickname),
		Content:       content,
		RelatedID:     commentID,
		TriggerUserID: likeUserID,
	})
}

// CreateVideoLikedNotification 视频被点赞
func (s *NotificationService) CreateVideoLikedNotification(videoID, likeUserID, toUserID uint64) {
	if likeUserID == toUserID {
		return
	}
	triggerUser, _ := s.userDao.GetByID(likeUserID)
	nickname := "某人"
	if triggerUser != nil {
		nickname = triggerUser.Nickname
		if nickname == "" {
			nickname = triggerUser.Username
		}
	}
	video, _ := s.videoDao.GetByID(videoID)
	videoTitle := ""
	if video != nil {
		videoTitle = video.Title
		if len(videoTitle) > 50 {
			videoTitle = videoTitle[:50] + "..."
		}
	}
	_ = s.notificationDao.Create(&model.Notification{
		UserID:        toUserID,
		Type:          model.NotificationTypeVideoLiked,
		Title:         fmt.Sprintf("%s 点赞了你的视频", nickname),
		Content:       videoTitle,
		RelatedID:     videoID,
		TriggerUserID: likeUserID,
	})
}

// ListNotifications 获取消息列表
func (s *NotificationService) ListNotifications(userID uint64, page, size int, unreadOnly bool) ([]model.Notification, int64, error) {
	return s.notificationDao.GetByUserID(userID, page, size, unreadOnly)
}

// GetUnreadCount 获取未读数量
func (s *NotificationService) GetUnreadCount(userID uint64) (int64, error) {
	return s.notificationDao.GetUnreadCount(userID)
}

// MarkAsRead 标记单条已读
func (s *NotificationService) MarkAsRead(id uint64) error {
	return s.notificationDao.MarkAsRead(id)
}

// MarkAllAsRead 全部已读
func (s *NotificationService) MarkAllAsRead(userID uint64) error {
	return s.notificationDao.MarkAllAsRead(userID)
}

// DeleteNotification 删除消息
func (s *NotificationService) DeleteNotification(id uint64) error {
	return s.notificationDao.Delete(id)
}
