package service

import (
	"time"

	"bilibili-backend/dao"
	"bilibili-backend/model"
)

type HistoryService struct {
	historyDao *dao.HistoryDao
}

func NewHistoryService(historyDao *dao.HistoryDao) *HistoryService {
	return &HistoryService{historyDao: historyDao}
}

// ListByUserID 查询用户播放历史（含视频信息）
func (s *HistoryService) ListByUserID(userID uint64, page, size int) ([]dao.HistoryVideo, int64, error) {
	return s.historyDao.ListWithVideo(userID, page, size)
}

// RecordHistory 记录播放历史（去重、限制 100 条、删除 30 天以上记录）
func (s *HistoryService) RecordHistory(userID, videoID uint64) error {
	// 1. 去重：先查询是否已有记录
	var existing model.VideoHistory
	err := s.historyDao.First(&existing, userID, videoID)
	if err == nil {
		// 已存在，更新 updated_at
		existing.UpdatedAt = time.Now()
		return s.historyDao.Save(&existing)
	}

	// 2. 新增记录
	h := &model.VideoHistory{
		UserID:  userID,
		VideoID: videoID,
	}
	if err := s.historyDao.Create(h); err != nil {
		return err
	}

	// 3. 清理超出 100 条的旧记录
	_ = s.historyDao.CleanupOverflow(userID, 100)

	// 4. 删除 30 天以上的记录
	_ = s.historyDao.CleanupOld(userID, 30*24*time.Hour)

	return nil
}
