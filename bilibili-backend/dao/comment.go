package dao

import (
	"bilibili-backend/model"
	"gorm.io/gorm"
)

type CommentDao struct {
	db *gorm.DB
}

func NewCommentDao(db *gorm.DB) *CommentDao {
	return &CommentDao{db: db}
}

func (d *CommentDao) Create(comment *model.Comment) error {
	return d.db.Create(comment).Error
}

func (d *CommentDao) GetByID(id uint64) (*model.Comment, error) {
	var c model.Comment
	if err := d.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *CommentDao) GetByIDWithUser(id uint64) (*model.Comment, error) {
	var c model.Comment
	if err := d.db.Preload("User").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// GetTopLevelByVideoID 获取视频的一级评论列表
func (d *CommentDao) GetTopLevelByVideoID(videoID uint64, page, size int, sort string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := d.db.Model(&model.Comment{}).Where("video_id = ? AND parent_id = 0", videoID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	order := "created_at DESC"
	if sort == "hot" {
		order = "like_count DESC, created_at DESC"
	}

	if err := query.Preload("User").Order(order).Limit(size).Offset(offset).Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// GetRepliesByRootID 获取某条一级评论的所有二级回复
func (d *CommentDao) GetRepliesByRootID(rootID uint64) ([]model.Comment, error) {
	var replies []model.Comment
	if err := d.db.Where("root_id = ?", rootID).Preload("User").Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}

// GetReplyCountByRootID 获取某条一级评论的回复数量
func (d *CommentDao) GetReplyCountByRootID(rootID uint64) (int64, error) {
	var count int64
	if err := d.db.Model(&model.Comment{}).Where("root_id = ?", rootID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *CommentDao) Delete(id uint64) error {
	return d.db.Delete(&model.Comment{}, id).Error
}

// IncrementLikeCount 点赞数 +/- delta
func (d *CommentDao) IncrementLikeCount(id uint64, delta int) error {
	return d.db.Model(&model.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

// UpdateVideoCommentCount 更新视频的 comment_count
func (d *CommentDao) UpdateVideoCommentCount(videoID uint64, count int) error {
	return d.db.Model(&model.Video{}).Where("id = ?", videoID).UpdateColumn("comment_count", count).Error
}
