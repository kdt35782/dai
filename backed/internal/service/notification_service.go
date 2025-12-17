package service

import (
	"sm-medical/internal/repository"
	"time"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		repo: repository.NewNotificationRepository(),
	}
}

// GetList 获取通知列表
func (s *NotificationService) GetList(userID int64, page, pageSize int, notifType string) ([]map[string]interface{}, int64, error) {
	notifications, total, err := s.repo.FindByUserID(userID, page, pageSize, notifType)
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, n := range notifications {
		item := map[string]interface{}{
			"notificationId": n.ID,
			"type":           n.NotificationType,
			"title":          n.Title,
			"content":        n.Content,
			"isRead":         n.IsRead,
			"createdAt":      n.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if n.RelatedID != nil {
			item["relatedId"] = *n.RelatedID
		}

		result = append(result, item)
	}

	return result, total, nil
}

// GetUnreadCount 获取未读数量
func (s *NotificationService) GetUnreadCount(userID int64) (map[string]interface{}, error) {
	total, system, consultation, err := s.repo.CountUnread(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalUnread":        total,
		"systemUnread":       system,
		"consultationUnread": consultation,
	}, nil
}

// MarkRead 标记已读
func (s *NotificationService) MarkRead(userID int64, notificationIDs []int64) error {
	now := time.Now()
	return s.repo.MarkAsRead(userID, notificationIDs, &now)
}
