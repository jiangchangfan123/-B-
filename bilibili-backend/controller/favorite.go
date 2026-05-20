package controller

import (
	"strconv"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type FavoriteController struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteController(favoriteService *service.FavoriteService) *FavoriteController {
	return &FavoriteController{favoriteService: favoriteService}
}

// ToggleFavorite POST /api/v1/videos/:id/favorite
func (ctrl *FavoriteController) ToggleFavorite(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	favorited, err := ctrl.favoriteService.ToggleFavorite(uid, videoID)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"favorited": favorited,
	})
}

// FavoriteStatus GET /api/v1/videos/:id/favorite/status
func (ctrl *FavoriteController) FavoriteStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	idStr := c.Param("id")
	videoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	favorited, _ := ctrl.favoriteService.IsFavorited(uid, videoID)
	utils.OK(c, gin.H{
		"favorited": favorited,
	})
}

// ListFavorites GET /api/v1/users/me/favorites
func (ctrl *FavoriteController) ListFavorites(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := ctrl.favoriteService.ListFavorites(uid, page, size)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
