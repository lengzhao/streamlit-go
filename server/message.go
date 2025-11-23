package server

// Message 消息结构
type Message struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ComponentEventData 组件事件数据
type ComponentEventData struct {
	ComponentID string `json:"componentId"`
	EventType   string `json:"eventType"`
	Value       string `json:"value"`
}

// MessageType 消息类型常量
const (
	MessageTypePing           = "ping"
	MessageTypePong           = "pong"
	MessageTypeUIUpdate       = "ui_update"
	MessageTypePartialUpdate  = "partial_update"
	MessageTypeError          = "error"
	MessageTypeComponentEvent = "component_event"
)
