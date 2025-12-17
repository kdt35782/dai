package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
	"time"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

// FindByUserID 根据用户ID查询
func (r *NotificationRepository) FindByUserID(userID int64, page, pageSize int, notifType string) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := database.GetDB().Model(&model.Notification{}).Where("user_id = ?", userID)

	if notifType != "" {
		query = query.Where("notification_type = ?", notifType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notifications).Error

	return notifications, total, err
}

// CountUnread 统计未读数量
func (r *NotificationRepository) CountUnread(userID int64) (int64, int64, int64, error) {
	var total, system, consultation int64

	database.GetDB().Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&total)
	database.GetDB().Model(&model.Notification{}).Where("user_id = ? AND is_read = ? AND notification_type = ?", userID, false, "system").Count(&system)
	database.GetDB().Model(&model.Notification{}).Where("user_id = ? AND is_read = ? AND notification_type = ?", userID, false, "consultation").Count(&consultation)

	return total, system, consultation, nil
}

// MarkAsRead 标记为已读
func (r *NotificationRepository) MarkAsRead(userID int64, notificationIDs []int64, readAt *time.Time) error {
	return database.GetDB().Model(&model.Notification{}).
		Where("user_id = ? AND id IN ?", userID, notificationIDs).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": readAt,
		}).Error
}
