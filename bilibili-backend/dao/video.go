package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type VideoDao struct {
	db *gorm.DB
}

func NewVideoDao(db *gorm.DB) *VideoDao {
	return &VideoDao{db: db}
}

func (d *VideoDao) ListByUserID(userID uint64, page, size int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64

	offset := (page - 1) * size
	if err := d.db.Model(&model.Video{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(size).Offset(offset).Find(&videos).Error; err != nil {
		return nil, 0, err
	}
	return videos, total, nil
}

func (d *VideoDao) GetByID(id uint64) (*model.Video, error) {
	var v model.Video
	if err := d.db.First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}
