package service

import (
	"bilibili-backend/dao"
	"bilibili-backend/model"
)

type FavoriteService struct {
	favoriteDao *dao.FavoriteDao
}

func NewFavoriteService(favoriteDao *dao.FavoriteDao) *FavoriteService {
	return &FavoriteService{favoriteDao: favoriteDao}
}

// ToggleFavorite 收藏/取消收藏
func (s *FavoriteService) ToggleFavorite(userID, videoID uint64) (favorited bool, err error) {
	fav, err := s.favoriteDao.GetByUserAndVideo(userID, videoID)
	if err != nil {
		// 未收藏，执行收藏
		f := &model.VideoFavorite{
			UserID:  userID,
			VideoID: videoID,
		}
		if err := s.favoriteDao.Create(f); err != nil {
			return false, err
		}
		return true, nil
	}
	// 已收藏，取消
	_ = fav
	if err := s.favoriteDao.Delete(userID, videoID); err != nil {
		return true, err
	}
	return false, nil
}

// IsFavorited 是否已收藏
func (s *FavoriteService) IsFavorited(userID, videoID uint64) (bool, error) {
	_, err := s.favoriteDao.GetByUserAndVideo(userID, videoID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ListFavorites 收藏列表
func (s *FavoriteService) ListFavorites(userID uint64, page, size int) ([]dao.FavoriteVideo, int64, error) {
	return s.favoriteDao.ListWithVideo(userID, page, size)
}
