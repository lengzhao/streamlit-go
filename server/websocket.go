package server

import (
"context"
"encoding/json"
"log"
"net/http"
"sync"
"time"

"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	hub        *Hub          // Hub引用
	conn       *websocket.Conn // WebSocket连接
	send       chan []byte   // 发送消息通道
	sessionID  string        // 会话ID
	mutex      sync.Mutex    // 互斥锁
}

// Upgrader WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源
		return true
	},
}

// NewClient 创建新的客户端
func NewClient(hub *Hub, conn *websocket.Conn, sessionID string) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 256),
		sessionID: sessionID,
	}
}

// SessionID 获取会话ID
func (c *Client) SessionID() string {
	return c.sessionID
}

// ReadPump 读取来自WebSocket连接的消息
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	// 设置读取超时
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			// 处理消息
			c.handleMessage(message)
		}
	}
}

// WritePump 将消息写入WebSocket连接
func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(54 * time.Second) // 每54秒发送一次ping
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的其他消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send 发送消息到客户端
func (c *Client) Send(message []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	select {
	case c.send <- message:
	default:
		// 如果通道满了，关闭连接
		close(c.send)
		c.conn.Close()
	}
}

// Close 关闭客户端连接
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conn.Close()
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(message []byte) {
	log.Printf("Client %s received message: %s", c.sessionID, string(message))

	// 解析消息
	msg := &Message{}
	if err := json.Unmarshal(message, msg); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	// 设置会话ID
	msg.SessionID = c.sessionID

	// 根据消息类型处理
	switch msg.Type {
	case MessageTypePing:
		// 回复pong消息
		pongMsg := &Message{
			Type:      MessageTypePong,
			SessionID: c.sessionID,
			Timestamp: time.Now().Unix(),
		}
		data, _ := json.Marshal(pongMsg)
		c.Send(data)
	case MessageTypeComponentEvent:
		// 处理组件事件
		c.hub.HandleComponentEvent(msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// ServeWebSocket 处理WebSocket连接
func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 获取会话ID
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		log.Println("Missing sessionId parameter")
		http.Error(w, "Missing sessionId parameter", http.StatusBadRequest)
		return
	}

	// 升级到WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// 创建客户端
	client := NewClient(hub, conn, sessionID)
	client.hub.Register(client)

	// 创建上下文用于取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动读写协程
	go client.WritePump(ctx)
	go client.ReadPump(ctx)

	// 发送初始UI
	// 注意：这里需要获取初始UI内容并发送给客户端
}
