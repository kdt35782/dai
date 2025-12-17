package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	Hub            *ChatHub
	Conn           *websocket.Conn
	Send           chan []byte
	UserID         int64
	ConsultationID int64
}

// ChatHub WebSocket连接管理中心
type ChatHub struct {
	// 所有客户端连接 key: consultationID_userID
	Clients map[string]*Client

	// 按问诊ID分组的客户端 key: consultationID
	ConsultationClients map[int64][]*Client

	// 按用户ID索引的客户端 key: userID
	UserClients map[int64][]*Client

	// 注册请求
	Register chan *Client

	// 注销请求
	Unregister chan *Client

	// 广播消息
	Broadcast chan *BroadcastMessage

	// 互斥锁
	mu sync.RWMutex
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	ConsultationID int64       // 问诊ID
	TargetUserID   int64       // 目标用户ID(0表示发送给问诊的所有人)
	MessageType    string      // 消息类型: "chat", "status", "typing", "notification"
	Data           interface{} // 消息数据
}

// NewChatHub 创建WebSocket管理中心
func NewChatHub() *ChatHub {
	return &ChatHub{
		Clients:             make(map[string]*Client),
		ConsultationClients: make(map[int64][]*Client),
		UserClients:         make(map[int64][]*Client),
		Register:            make(chan *Client),
		Unregister:          make(chan *Client),
		Broadcast:           make(chan *BroadcastMessage, 256),
	}
}

// Run 运行WebSocket中心
func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *ChatHub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientKey := getClientKey(client.ConsultationID, client.UserID)
	h.Clients[clientKey] = client

	// 添加到问诊分组
	h.ConsultationClients[client.ConsultationID] = append(
		h.ConsultationClients[client.ConsultationID],
		client,
	)

	// 添加到用户索引
	h.UserClients[client.UserID] = append(h.UserClients[client.UserID], client)

	log.Printf("[WebSocket] 客户端已连接 - 用户ID: %d, 问诊ID: %d, 当前连接数: %d",
		client.UserID, client.ConsultationID, len(h.Clients))

	// 发送连接成功消息
	h.sendToClient(client, map[string]interface{}{
		"type": "connected",
		"data": map[string]interface{}{
			"userId":         client.UserID,
			"consultationId": client.ConsultationID,
			"timestamp":      time.Now().Unix(),
		},
	})

	// 通知其他用户该用户上线
	h.broadcastMessage(&BroadcastMessage{
		ConsultationID: client.ConsultationID,
		TargetUserID:   0,
		MessageType:    "status",
		Data: map[string]interface{}{
			"userId": client.UserID,
			"status": "online",
		},
	})
}

// unregisterClient 注销客户端
func (h *ChatHub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientKey := getClientKey(client.ConsultationID, client.UserID)
	if _, ok := h.Clients[clientKey]; ok {
		delete(h.Clients, clientKey)
		close(client.Send)

		// 从问诊分组中移除
		clients := h.ConsultationClients[client.ConsultationID]
		for i, c := range clients {
			if c == client {
				h.ConsultationClients[client.ConsultationID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}

		// 从用户索引中移除
		userClients := h.UserClients[client.UserID]
		for i, c := range userClients {
			if c == client {
				h.UserClients[client.UserID] = append(userClients[:i], userClients[i+1:]...)
				break
			}
		}

		log.Printf("[WebSocket] 客户端已断开 - 用户ID: %d, 问诊ID: %d, 剩余连接数: %d",
			client.UserID, client.ConsultationID, len(h.Clients))

		// 通知其他用户该用户下线
		h.broadcastMessage(&BroadcastMessage{
			ConsultationID: client.ConsultationID,
			TargetUserID:   0,
			MessageType:    "status",
			Data: map[string]interface{}{
				"userId": client.UserID,
				"status": "offline",
			},
		})
	}
}

// broadcastMessage 广播消息
func (h *ChatHub) broadcastMessage(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	message := map[string]interface{}{
		"type":      msg.MessageType,
		"data":      msg.Data,
		"timestamp": time.Now().Unix(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("[WebSocket] 消息序列化失败: %v", err)
		return
	}

	// 如果指定了目标用户,只发送给该用户
	if msg.TargetUserID > 0 {
		clientKey := getClientKey(msg.ConsultationID, msg.TargetUserID)
		if client, ok := h.Clients[clientKey]; ok {
			select {
			case client.Send <- messageBytes:
			default:
				log.Printf("[WebSocket] 发送失败,通道已满 - 用户ID: %d", msg.TargetUserID)
			}
		}
		return
	}

	// 否则发送给问诊的所有在线用户
	clients := h.ConsultationClients[msg.ConsultationID]
	for _, client := range clients {
		select {
		case client.Send <- messageBytes:
		default:
			log.Printf("[WebSocket] 发送失败,通道已满 - 用户ID: %d", client.UserID)
		}
	}
}

// sendToClient 发送消息给指定客户端
func (h *ChatHub) sendToClient(client *Client, message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("[WebSocket] 消息序列化失败: %v", err)
		return
	}

	select {
	case client.Send <- messageBytes:
	default:
		log.Printf("[WebSocket] 发送失败,通道已满 - 用户ID: %d", client.UserID)
	}
}

// IsUserOnline 检查用户是否在线
func (h *ChatHub) IsUserOnline(consultationID, userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clientKey := getClientKey(consultationID, userID)
	_, ok := h.Clients[clientKey]
	return ok
}

// GetOnlineUsers 获取问诊中的在线用户列表
func (h *ChatHub) GetOnlineUsers(consultationID int64) []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := h.ConsultationClients[consultationID]
	userIDs := make([]int64, 0, len(clients))
	userMap := make(map[int64]bool)

	for _, client := range clients {
		if !userMap[client.UserID] {
			userIDs = append(userIDs, client.UserID)
			userMap[client.UserID] = true
		}
	}

	return userIDs
}

// SendToUser 发送消息给指定用户
func (h *ChatHub) SendToUser(consultationID, userID int64, messageType string, data interface{}) {
	h.Broadcast <- &BroadcastMessage{
		ConsultationID: consultationID,
		TargetUserID:   userID,
		MessageType:    messageType,
		Data:           data,
	}
}

// SendToConsultation 发送消息给问诊的所有用户
func (h *ChatHub) SendToConsultation(consultationID int64, messageType string, data interface{}) {
	h.Broadcast <- &BroadcastMessage{
		ConsultationID: consultationID,
		TargetUserID:   0,
		MessageType:    messageType,
		Data:           data,
	}
}

// getClientKey 生成客户端唯一键
func getClientKey(consultationID, userID int64) string {
	return fmt.Sprintf("%d_%d", consultationID, userID)
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] 读取错误: %v", err)
			}
			break
		}

		// 处理客户端消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[WebSocket] 消息解析失败: %v", err)
			continue
		}

		// 根据消息类型处理
		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}

		switch msgType {
		case "ping":
			// 心跳响应
			c.Hub.sendToClient(c, map[string]interface{}{
				"type": "pong",
			})

		case "typing":
			// 转发正在输入状态
			c.Hub.SendToConsultation(c.ConsultationID, "typing", map[string]interface{}{
				"userId": c.UserID,
				"typing": msg["data"],
			})
		}
	}
}

// WritePump 写入消息到客户端
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送队列中的其他消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
