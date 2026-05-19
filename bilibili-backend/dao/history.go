package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type HistoryDao struct {
	db *gorm.DB
}

func NewHistoryDao(db *gorm.DB) *HistoryDao {
	return &HistoryDao{db: db}
}

func (d *HistoryDao) ListByUserID(userID uint64, page, size int) ([]model.VideoHistory, int64, error) {
	var histories []model.VideoHistory
	var total int64

	offset := (page - 1) * size
	if err := d.db.Model(&model.VideoHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := d.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(size).Offset(offset).Find(&histories).Error; err != nil {
		return nil, 0, err
	}
	return histories, total, nil
}

func (d *HistoryDao) Create(h *model.VideoHistory) error {
	return d.db.Create(h).Error
}
