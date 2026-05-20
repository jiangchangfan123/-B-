package service

import (
	"encoding/json"
	"log"

	"bilibili-backend/dao"
	"bilibili-backend/model"
	"bilibili-backend/utils"

	"gorm.io/gorm"
)

// DanmakuService 弹幕业务逻辑
type DanmakuService struct {
	danmakuDao *dao.DanmakuDAO
	db         *gorm.DB
}

// NewDanmakuService 创建弹幕服务
func NewDanmakuService(danmakuDao *dao.DanmakuDAO, db *gorm.DB) *DanmakuService {
	return &DanmakuService{
		danmakuDao: danmakuDao,
		db:         db,
	}
}

// CreateDanmaku 创建弹幕并广播，返回创建的弹幕数据
func (s *DanmakuService) CreateDanmaku(videoID, userID uint64, content string, timePoint int, color string, dType int8) (*model.Danmaku, error) {
	danmaku := &model.Danmaku{
		VideoID:   videoID,
		UserID:    userID,
		Content:   content,
		TimePoint: timePoint,
		Color:     color,
		Type:      dType,
	}

	if err := s.danmakuDao.Create(s.db, danmaku); err != nil {
		return nil, err
	}

	// 广播给同视频的所有 WebSocket 连接
	go s.broadcastDanmaku(videoID, danmaku)

	return danmaku, nil
}

// GetDanmakuByTimeRange 获取某视频某时间段的弹幕
func (s *DanmakuService) GetDanmakuByTimeRange(videoID uint64, start, end int) ([]model.Danmaku, error) {
	return s.danmakuDao.GetByVideoIDTimeRange(s.db, videoID, start, end)
}

// broadcastDanmaku 将弹幕广播给同视频的所有连接
func (s *DanmakuService) broadcastDanmaku(videoID uint64, danmaku *model.Danmaku) {
	msg, err := json.Marshal(map[string]interface{}{
		"type": "danmaku",
		"data": danmaku,
	})
	if err != nil {
		log.Printf("[Danmaku] JSON 序列化失败: %v", err)
		return
	}

	utils.DanmakuHubInstance.Broadcast(int64(videoID), msg)
	log.Printf("[Danmaku] 已广播到 video_id=%d", videoID)
}
