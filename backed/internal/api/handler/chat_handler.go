package handler

import (
	"log"
	"net/http"
	"sm-medical/internal/service"
	"sm-medical/internal/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var (
	chatService *service.ChatService
	chatHub     *websocket.ChatHub
	upgrader    = ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 生产环境需要验证Origin
		},
	}
)

// InitChatHandler 初始化聊天处理器
func InitChatHandler() {
	chatService = service.NewChatService()
	chatHub = websocket.NewChatHub()
	
	// 启动WebSocket中心
	go chatHub.Run()
	
	log.Println("[ChatHandler] 聊天服务已初始化")
}

// WebSocketConnect WebSocket连接处理
func WebSocketConnect(c *gin.Context) {
	// 获取参数
	userIDStr := c.Query("userId")
	consultationIDStr := c.Query("consultationId")

	if userIDStr == "" || consultationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要参数",
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}

	consultationID, err := strconv.ParseInt(consultationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "问诊ID格式错误",
		})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] 升级失败: %v", err)
		return
	}

	// 创建客户端
	client := &websocket.Client{
		Hub:            chatHub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		UserID:         userID,
		ConsultationID: consultationID,
	}

	// 注册客户端
	chatHub.Register <- client

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	var req service.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层发送消息
	message, err := chatService.SendMessage(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	// 通过WebSocket推送给接收者
	chatHub.SendToUser(message.ConsultationID, message.ReceiverID, "chat", message)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送成功",
		"data":    message,
	})
}

// GetMessageList 获取消息列表
func GetMessageList(c *gin.Context) {
	var req service.GetMessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := chatService.GetMessageList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    result,
	})
}

// GetUnreadCount 获取未读消息数量
func GetUnreadCount(c *gin.Context) {
	var req service.GetUnreadCountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := chatService.GetUnreadCount(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    result,
	})
}

// MarkAsRead 标记消息已读
func MarkAsRead(c *gin.Context) {
	var req service.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := chatService.MarkAsRead(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "标记成功",
	})
}

// GetOnlineStatus 获取在线状态
func GetOnlineStatus(c *gin.Context) {
	consultationIDStr := c.Query("consultationId")
	if consultationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少问诊ID",
		})
		return
	}

	consultationID, err := strconv.ParseInt(consultationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "问诊ID格式错误",
		})
		return
	}

	onlineUsers := chatHub.GetOnlineUsers(consultationID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"consultationId": consultationID,
			"onlineUsers":    onlineUsers,
		},
	})
}

// SendTypingStatus 发送正在输入状态
func SendTypingStatus(c *gin.Context) {
	var req struct {
		ConsultationID int64 `json:"consultationId" binding:"required"`
		UserID         int64 `json:"userId" binding:"required"`
		Typing         bool  `json:"typing"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 广播正在输入状态
	chatHub.SendToConsultation(req.ConsultationID, "typing", gin.H{
		"userId": req.UserID,
		"typing": req.Typing,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送成功",
	})
}
