package repository

import (
	"sm-medical/internal/model"
	"sm-medical/pkg/database"
	"time"
)

type ChatRepository struct{}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

// CreateMessage 创建消息
func (r *ChatRepository) CreateMessage(message *model.ChatMessage) error {
	return database.DB.Create(message).Error
}

// GetMessagesByConsultationID 获取问诊的所有消息
func (r *ChatRepository) GetMessagesByConsultationID(consultationID int64, page, pageSize int) ([]model.ChatMessage, int64, error) {
	var messages []model.ChatMessage
	var total int64

	db := database.DB.Model(&model.ChatMessage{}).
		Where("consultation_id = ? AND is_deleted = 0", consultationID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// GetUnreadMessages 获取未读消息
func (r *ChatRepository) GetUnreadMessages(userID, consultationID int64) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	err := database.DB.Where("consultation_id = ? AND receiver_id = ? AND is_read = 0 AND is_deleted = 0", consultationID, userID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// MarkAsRead 标记消息已读 (批量)
func (r *ChatRepository) MarkAsRead(messageIDs []int64) error {
	now := time.Now()
	return database.DB.Model(&model.ChatMessage{}).
		Where("id IN ?", messageIDs).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkMessageAsRead 标记单个消息已读
func (r *ChatRepository) MarkMessageAsRead(messageID int64) error {
	now := time.Now()
	return database.DB.Model(&model.ChatMessage{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead 标记问诊的所有消息已读
func (r *ChatRepository) MarkAllAsRead(userID, consultationID int64) error {
	now := time.Now()
	return database.DB.Model(&model.ChatMessage{}).
		Where("consultation_id = ? AND receiver_id = ? AND is_read = 0", consultationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// GetUnreadCount 获取未读消息数量
func (r *ChatRepository) GetUnreadCount(userID, consultationID int64) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ChatMessage{}).
		Where("consultation_id = ? AND receiver_id = ? AND is_read = 0 AND is_deleted = 0", consultationID, userID).
		Count(&count).Error
	return count, err
}

// UpdateUnreadCount 更新未读消息统计
func (r *ChatRepository) UpdateUnreadCount(userID, consultationID int64, count int, lastMessageID int64, lastMessageTime time.Time) error {
	// 使用ON DUPLICATE KEY UPDATE
	return database.DB.Exec(`
		INSERT INTO SM_chat_unread_count (user_id, consultation_id, unread_count, last_message_id, last_message_time)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			unread_count = VALUES(unread_count),
			last_message_id = VALUES(last_message_id),
			last_message_time = VALUES(last_message_time)
	`, userID, consultationID, count, lastMessageID, lastMessageTime).Error
}

// GetUnreadCountByUser 获取用户所有问诊的未读消息统计
func (r *ChatRepository) GetUnreadCountByUser(userID int64) ([]model.ChatUnreadCount, error) {
	var counts []model.ChatUnreadCount
	err := database.DB.Where("user_id = ? AND unread_count > 0", userID).
		Order("last_message_time DESC").
		Find(&counts).Error
	return counts, err
}

// GetUnreadCountList 获取用户所有问诊的未读消息统计 (别名)
func (r *ChatRepository) GetUnreadCountList(userID int64) ([]model.ChatUnreadCount, error) {
	return r.GetUnreadCountByUser(userID)
}

// GetMessageByID 根据ID获取消息
func (r *ChatRepository) GetMessageByID(messageID int64) (*model.ChatMessage, error) {
	var message model.ChatMessage
	err := database.DB.Where("id = ? AND is_deleted = 0", messageID).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// DeleteMessage 软删除消息
func (r *ChatRepository) DeleteMessage(messageID int64) error {
	return database.DB.Model(&model.ChatMessage{}).
		Where("id = ?", messageID).
		Update("is_deleted", true).Error
}

// GetLastMessage 获取问诊的最后一条消息
func (r *ChatRepository) GetLastMessage(consultationID int64) (*model.ChatMessage, error) {
	var message model.ChatMessage
	err := database.DB.Where("consultation_id = ? AND is_deleted = 0", consultationID).
		Order("created_at DESC").
		First(&message).Error
	return &message, err
}
