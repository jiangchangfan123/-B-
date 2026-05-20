package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type NotificationDao struct {
	db *gorm.DB
}

func NewNotificationDao(db *gorm.DB) *NotificationDao {
	return &NotificationDao{db: db}
}

func (d *NotificationDao) Create(n *model.Notification) error {
	return d.db.Create(n).Error
}

func (d *NotificationDao) GetByUserID(userID uint64, page, size int, unreadOnly bool) ([]model.Notification, int64, error) {
	query := d.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Notification
	err := query.Order("created_at DESC").Offset((page - 1) * size).Limit(size).
		Preload("TriggerUser").Find(&list).Error
	return list, total, err
}

func (d *NotificationDao) GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	err := d.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (d *NotificationDao) MarkAsRead(id uint64) error {
	return d.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (d *NotificationDao) MarkAllAsRead(userID uint64) error {
	return d.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (d *NotificationDao) Delete(id uint64) error {
	return d.db.Delete(&model.Notification{}, id).Error
}
