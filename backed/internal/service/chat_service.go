package service

import (
	"errors"
	"fmt"
	"log"
	"sm-medical/internal/crypto"
	"sm-medical/internal/model"
	"sm-medical/internal/repository"
	"time"
)

// ChatService 聊天服务
type ChatService struct {
	chatRepo         *repository.ChatRepository
	consultationRepo *repository.ConsultationRepository
	userRepo         *repository.UserRepository
}

// NewChatService 创建聊天服务实例
func NewChatService() *ChatService {
	return &ChatService{
		chatRepo:         repository.NewChatRepository(),
		consultationRepo: repository.NewConsultationRepository(),
		userRepo:         repository.NewUserRepository(),
	}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ConsultationID int64  `json:"consultationId" binding:"required"`
	SenderID       int64  `json:"senderId" binding:"required"`
	MessageType    int    `json:"messageType" binding:"required"` // 1:文本 2:图片 3:语音 4:处方 5:系统
	Content        string `json:"content"`
	FileURL        string `json:"fileUrl"`
	FileSize       int    `json:"fileSize"`
	Duration       int    `json:"duration"`
	ExtraData      string `json:"extraData"`
}

// SendMessage 发送消息
func (s *ChatService) SendMessage(req *SendMessageRequest) (*model.ChatMessage, error) {
	// 1. 验证问诊存在性
	consultation, err := s.consultationRepo.FindByID(req.ConsultationID)
	if err != nil {
		return nil, errors.New("问诊不存在")
	}

	// 2. 确定接收者
	var receiverID int64
	if consultation.DoctorID != nil && req.SenderID == *consultation.DoctorID {
		receiverID = consultation.PatientID
	} else if req.SenderID == consultation.PatientID {
		if consultation.DoctorID != nil {
			receiverID = *consultation.DoctorID
		} else {
			return nil, errors.New("问诊尚未分配医生")
		}
	} else {
		return nil, errors.New("无权限发送消息")
	}

	// 3. 生成消息编号
	messageNo := generateMessageNo()

	// 4. 加密消息内容
	var encryptedContent string
	if req.Content != "" {
		encrypted, err := crypto.SM4Encrypt(req.Content)
		if err != nil {
			log.Printf("[ChatService] SM4加密失败: %v", err)
			return nil, errors.New("消息加密失败")
		}
		encryptedContent = encrypted
	}

	// 5. 创建消息
	message := &model.ChatMessage{
		MessageNo:      messageNo,
		ConsultationID: req.ConsultationID,
		SenderID:       req.SenderID,
		ReceiverID:     receiverID,
		MessageType:    req.MessageType,
		Content:        encryptedContent,
		FileURL:        req.FileURL,
		FileSize:       req.FileSize,
		Duration:       req.Duration,
		ExtraData:      req.ExtraData,
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	if err := s.chatRepo.CreateMessage(message); err != nil {
		log.Printf("[ChatService] 创建消息失败: %v", err)
		return nil, err
	}

	// 6. 更新接收者的未读消息统计
	unreadCount, _ := s.chatRepo.GetUnreadCount(receiverID, req.ConsultationID)
	err = s.chatRepo.UpdateUnreadCount(receiverID, req.ConsultationID, int(unreadCount+1), message.ID, message.CreatedAt)
	if err != nil {
		log.Printf("[ChatService] 更新未读统计失败: %v", err)
	}

	// 7. 填充发送者信息
	sender, _ := s.userRepo.FindByID(req.SenderID)
	if sender != nil {
		message.SenderName = sender.Username // 使用 Username
		message.SenderAvatar = sender.Avatar
		if sender.Role == "patient" {
			message.SenderRole = "patient"
		} else if sender.Role == "doctor" {
			message.SenderRole = "doctor"
		}
	}

	log.Printf("[ChatService] 消息发送成功 - 消息ID: %d, 问诊ID: %d, 发送者: %d, 接收者: %d", 
		message.ID, req.ConsultationID, req.SenderID, receiverID)

	return message, nil
}

// GetMessageListRequest 获取消息列表请求
type GetMessageListRequest struct {
	ConsultationID int64 `form:"consultationId" binding:"required"`
	UserID         int64 `form:"userId" binding:"required"`
	Page           int   `form:"page"`
	PageSize       int   `form:"pageSize"`
}

// GetMessageListResponse 消息列表响应
type GetMessageListResponse struct {
	Messages []model.ChatMessage `json:"messages"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

// GetMessageList 获取消息列表
func (s *ChatService) GetMessageList(req *GetMessageListRequest) (*GetMessageListResponse, error) {
	// 1. 验证权限
	consultation, err := s.consultationRepo.FindByID(req.ConsultationID)
	if err != nil {
		return nil, errors.New("问诊不存在")
	}

	if req.UserID != consultation.PatientID && (consultation.DoctorID == nil || req.UserID != *consultation.DoctorID) {
		return nil, errors.New("无权限查看消息")
	}

	// 2. 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}

	// 3. 获取消息列表
	messages, total, err := s.chatRepo.GetMessagesByConsultationID(req.ConsultationID, req.Page, req.PageSize)
	if err != nil {
		log.Printf("[ChatService] 获取消息列表失败: %v", err)
		return nil, err
	}

	// 4. 解密消息内容并填充发送者信息
	for i := range messages {
		// 解密内容
		if messages[i].Content != "" {
			decrypted, err := crypto.SM4Decrypt(messages[i].Content)
			if err != nil {
				log.Printf("[ChatService] 消息解密失败 ID=%d: %v", messages[i].ID, err)
				messages[i].Content = "[解密失败]"
			} else {
				messages[i].Content = decrypted
			}
		}

		// 填充发送者信息
		sender, _ := s.userRepo.FindByID(messages[i].SenderID)
		if sender != nil {
			messages[i].SenderName = sender.Username // 使用 Username
			messages[i].SenderAvatar = sender.Avatar
			if sender.Role == "patient" {
				messages[i].SenderRole = "patient"
			} else if sender.Role == "doctor" {
				messages[i].SenderRole = "doctor"
			}
		}
	}

	// 5. 标记所有消息为已读
	if err := s.chatRepo.MarkAllAsRead(req.UserID, req.ConsultationID); err != nil {
		log.Printf("[ChatService] 标记已读失败: %v", err)
	}

	// 6. 清空未读统计
	s.chatRepo.UpdateUnreadCount(req.UserID, req.ConsultationID, 0, 0, time.Time{})

	return &GetMessageListResponse{
		Messages: messages,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetUnreadCountRequest 获取未读数量请求
type GetUnreadCountRequest struct {
	UserID         int64 `form:"userId" binding:"required"`
	ConsultationID int64 `form:"consultationId"`
}

// UnreadCountItem 未读数量项
type UnreadCountItem struct {
	ConsultationID  int64      `json:"consultationId"`
	UnreadCount     int        `json:"unreadCount"`
	LastMessageID   *int64     `json:"lastMessageId"`
	LastMessageTime *time.Time `json:"lastMessageTime"`
}

// GetUnreadCount 获取未读消息数量
func (s *ChatService) GetUnreadCount(req *GetUnreadCountRequest) (interface{}, error) {
	// 如果指定了问诊ID,返回单个问诊的未读数
	if req.ConsultationID > 0 {
		count, err := s.chatRepo.GetUnreadCount(req.UserID, req.ConsultationID)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"consultationId": req.ConsultationID,
			"unreadCount":    count,
		}, nil
	}

	// 否则返回所有问诊的未读统计
	unreadList, err := s.chatRepo.GetUnreadCountList(req.UserID)
	if err != nil {
		return nil, err
	}

	var result []UnreadCountItem
	for _, item := range unreadList {
		result = append(result, UnreadCountItem{
			ConsultationID:  item.ConsultationID,
			UnreadCount:     item.UnreadCount,
			LastMessageID:   item.LastMessageID,
			LastMessageTime: item.LastMessageTime,
		})
	}

	return result, nil
}

// MarkAsReadRequest 标记已读请求
type MarkAsReadRequest struct {
	MessageID int64 `json:"messageId" binding:"required"`
	UserID    int64 `json:"userId" binding:"required"`
}

// MarkAsRead 标记消息已读
func (s *ChatService) MarkAsRead(req *MarkAsReadRequest) error {
	// 获取消息
	message, err := s.chatRepo.GetMessageByID(req.MessageID)
	if err != nil {
		return errors.New("消息不存在")
	}

	// 验证权限(只有接收者才能标记已读)
	if message.ReceiverID != req.UserID {
		return errors.New("无权限标记该消息")
	}

	// 标记已读
	if err := s.chatRepo.MarkMessageAsRead(req.MessageID); err != nil {
		return err
	}

	// 更新未读统计
	unreadCount, _ := s.chatRepo.GetUnreadCount(req.UserID, message.ConsultationID)
	s.chatRepo.UpdateUnreadCount(req.UserID, message.ConsultationID, int(unreadCount), 0, time.Time{})

	return nil
}

// SendSystemMessage 发送系统消息(内部方法)
func (s *ChatService) SendSystemMessage(consultationID int64, receiverID int64, content string) error {
	messageNo := generateMessageNo()

	// 加密内容
	encryptedContent, err := crypto.SM4Encrypt(content)
	if err != nil {
		log.Printf("[ChatService] 系统消息加密失败: %v", err)
		return err
	}

	message := &model.ChatMessage{
		MessageNo:      messageNo,
		ConsultationID: consultationID,
		SenderID:       0, // 系统消息发送者ID为0
		ReceiverID:     receiverID,
		MessageType:    5, // 系统消息
		Content:        encryptedContent,
		IsRead:         false,
		CreatedAt:      time.Now(),
	}

	if err := s.chatRepo.CreateMessage(message); err != nil {
		log.Printf("[ChatService] 创建系统消息失败: %v", err)
		return err
	}

	// 更新未读统计
	unreadCount, _ := s.chatRepo.GetUnreadCount(receiverID, consultationID)
	s.chatRepo.UpdateUnreadCount(receiverID, consultationID, int(unreadCount+1), message.ID, message.CreatedAt)

	log.Printf("[ChatService] 系统消息发送成功 - 问诊 ID: %d, 接收者: %d", consultationID, receiverID)

	return nil
}

// generateMessageNo 生成消息编号
func generateMessageNo() string {
	return fmt.Sprintf("MSG%d%06d", time.Now().Unix(), time.Now().Nanosecond()%1000000)
}
