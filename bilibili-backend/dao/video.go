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

func (d *VideoDao) Create(video *model.Video) error {
	return d.db.Create(video).Error
}

func (d *VideoDao) GetByID(id uint64) (*model.Video, error) {
	var v model.Video
	if err := d.db.First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (d *VideoDao) Delete(id uint64) error {
	return d.db.Delete(&model.Video{}, id).Error
}

// ListByUserID 查询用户的视频列表
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

// ListPublished 查询已发布视频列表（带分页、分类、排序）
func (d *VideoDao) ListPublished(category string, sort string, page, size int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64

	query := d.db.Model(&model.Video{}).Where("status = ? AND transcode_status = ?", 1, 2)
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	order := "created_at DESC"
	if sort == "hot" {
		order = "view_count DESC"
	}

	if err := query.Order(order).Limit(size).Offset(offset).Find(&videos).Error; err != nil {
		return nil, 0, err
	}
	return videos, total, nil
}

// UpdateTranscodeStatus 更新转码状态和播放地址
func (d *VideoDao) UpdateTranscodeStatus(id uint64, transcodeStatus int, transcodedURL string, status int) error {
	updates := map[string]interface{}{
		"transcode_status": transcodeStatus,
		"status":           status,
	}
	if transcodedURL != "" {
		updates["transcoded_url"] = transcodedURL
	}
	return d.db.Model(&model.Video{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementViewCount 播放量 +1
func (d *VideoDao) IncrementViewCount(id uint64) error {
	return d.db.Model(&model.Video{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateLikeCount 更新点赞数
func (d *VideoDao) UpdateLikeCount(id uint64, count int) error {
	return d.db.Model(&model.Video{}).Where("id = ?", id).UpdateColumn("like_count", count).Error
}

// Update 更新视频字段
func (d *VideoDao) Update(video *model.Video) error {
	return d.db.Save(video).Error
}
