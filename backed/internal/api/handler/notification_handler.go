package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/service"
	"sm-medical/pkg/utils"
	"strconv"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		service: service.NewNotificationService(),
	}
}

// GetList 获取通知列表
func (h *NotificationHandler) GetList(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	notifType := c.Query("type")

	list, total, err := h.service.GetList(userID, page, pageSize, notifType)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetUnreadCount 获取未读数量
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt64("userID")

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, count)
}

// MarkRead 标记已读
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := c.GetInt64("userID")

	var req struct {
		NotificationIDs []int64 `json:"notificationIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.service.MarkRead(userID, req.NotificationIDs); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "已标记为已读", nil)
}
