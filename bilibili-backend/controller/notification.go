package controller

import (
	"strconv"

	"bilibili-backend/service"
	"bilibili-backend/utils"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService *service.NotificationService
}

func NewNotificationController(notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{notificationService: notificationService}
}

func (ctrl *NotificationController) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	unreadOnly := c.Query("unread_only") == "true"

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	list, total, err := ctrl.notificationService.ListNotifications(uid, page, size, unreadOnly)
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

func (ctrl *NotificationController) UnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	count, err := ctrl.notificationService.GetUnreadCount(uid)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, gin.H{"count": count})
}

func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	if err := ctrl.notificationService.MarkAsRead(id); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, nil)
}

func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	if err := ctrl.notificationService.MarkAllAsRead(uid); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, nil)
}

func (ctrl *NotificationController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	if err := ctrl.notificationService.DeleteNotification(id); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}
	utils.OK(c, nil)
}
