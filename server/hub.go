package server

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/lengzhao/streamlit-go/state"
)

// Hub WebSocket连接池管理器
type Hub struct {
	// 客户端连接池（按会话ID分组）
	clients map[string]map[*Client]bool

	// 注册客户端通道
	register chan *Client

	// 注销客户端通道
	unregister chan *Client

	// 广播消息通道
	broadcast chan []byte

	// 互斥锁
	mu sync.RWMutex

	// 组件事件处理器
	eventHandler EventHandler

	// 状态管理器引用
	stateManager *state.Manager
}

// EventHandler 事件处理器接口
type EventHandler interface {
	HandleComponentEvent(session *state.Session, event *ComponentEventData)
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessionID := client.SessionID()
	if h.clients[sessionID] == nil {
		h.clients[sessionID] = make(map[*Client]bool)
	}
	h.clients[sessionID][client] = true

	log.Printf("Client registered for session: %s, total clients: %d", sessionID, h.ClientCount())
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessionID := client.SessionID()
	if clients, ok := h.clients[sessionID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			client.Close()

			// 如果该会话没有客户端了，删除会话
			if len(clients) == 0 {
				delete(h.clients, sessionID)
			}

			log.Printf("Client unregistered for session: %s, total clients: %d", sessionID, h.ClientCount())
		}
	}
}

// broadcastMessage 广播消息到所有客户端
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.clients {
		for client := range clients {
			client.Send(message)
		}
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	// log.Printf("Unregistering client for session: %s", client.SessionID())
	h.unregister <- client
}

// Broadcast 广播消息
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// SendToSession 发送消息到指定会话的所有客户端
func (h *Hub) SendToSession(sessionID string, message []byte) {
	log.Printf("Hub.SendToSession: sessionID=%s, message=%s", sessionID, string(message))
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[sessionID]; ok {
		log.Printf("Hub.SendToSession: found %d clients for session %s", len(clients), sessionID)
		count := 0
		for client := range clients {
			log.Printf("Hub.SendToSession: sending to client %d for session %s", count, sessionID)
			client.Send(message)
			count++
		}
		log.Printf("Hub.SendToSession: sent to %d clients for session %s", count, sessionID)
	} else {
		log.Printf("Hub.SendToSession: no clients found for session %s", sessionID)
	}
}

// ClientCount 返回当前连接的客户端总数
func (h *Hub) ClientCount() int {
	count := 0
	for _, clients := range h.clients {
		count += len(clients)
	}
	return count
}

// SessionCount 返回当前会话数量
func (h *Hub) SessionCount() int {
	return len(h.clients)
}

// SetEventHandler 设置事件处理器
func (h *Hub) SetEventHandler(handler EventHandler) {
	h.eventHandler = handler
}

// SetStateManager 设置状态管理器
func (h *Hub) SetStateManager(stateManager *state.Manager) {
	h.stateManager = stateManager
}

// HandleComponentEvent 处理组件事件
func (h *Hub) HandleComponentEvent(msg *Message) {
	if h.eventHandler == nil {
		log.Println("No event handler set")
		return
	}

	// 解析组件事件数据
	eventData := &ComponentEventData{}
	if msg.Data != nil {
		d, _ := json.Marshal(msg.Data)
		_ = json.Unmarshal(d, eventData)
	}

	// 获取会话对象
	session := h.stateManager.GetSession(msg.SessionID)

	// 首先尝试在会话中处理事件
	h.eventHandler.HandleComponentEvent(session, eventData)
}

// SendUIUpdate 发送UI更新到指定会话
func (h *Hub) SendUIUpdate(sessionID string, html string) error {
	msg := Message{
		Type:      MessageTypeUIUpdate,
		SessionID: sessionID,
		Data: map[string]interface{}{
			"html": html,
		},
		Timestamp: 0, // 将在发送时设置
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.SendToSession(sessionID, data)
	return nil
}

// SendPartialUpdate 发送局部更新到指定会话
func (h *Hub) SendPartialUpdate(sessionID string, componentID string, html string) error {
	log.Printf("Hub.SendPartialUpdate: sessionID=%s, componentID=%s", sessionID, componentID)
	msg := Message{
		Type:      MessageTypePartialUpdate,
		SessionID: sessionID,
		Data: map[string]interface{}{
			"componentId": componentID,
			"html":        html,
		},
		Timestamp: 0,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.SendToSession(sessionID, data)
	return nil
}

// SendError 发送错误消息到指定会话
func (h *Hub) SendError(sessionID string, errorMsg string, errorCode string) error {
	msg := Message{
		Type:      MessageTypeError,
		SessionID: sessionID,
		Data: map[string]interface{}{
			"message": errorMsg,
			"code":    errorCode,
		},
		Timestamp: 0,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.SendToSession(sessionID, data)
	return nil
}

// SendAddWidget 发送添加组件消息到指定会话
func (h *Hub) SendAddWidget(sessionID string, componentID string, html string) error {
	log.Printf("Hub.SendAddWidget: sessionID=%s, componentID=%s", sessionID, componentID)
	msg := Message{
		Type:      "add_widget",
		SessionID: sessionID,
		Data: map[string]interface{}{
			"componentId": componentID,
			"html":        html,
		},
		Timestamp: 0,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.SendToSession(sessionID, data)
	return nil
}
