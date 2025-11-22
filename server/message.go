package server

// MessageType 消息类型
type MessageType string

const (
MessageTypePing          MessageType = "ping"
MessageTypePong          MessageType = "pong"
MessageTypeComponentEvent MessageType = "component_event"
MessageTypeUIUpdate      MessageType = "ui_update"
MessageTypePartialUpdate MessageType = "partial_update"
MessageTypeError         MessageType = "error"
)

// Message 消息结构
type Message struct {
	Type      MessageType            `json:"type"`
	SessionID string                 `json:"sessionId"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp int64                  `json:"timestamp"`
}

// ComponentEventData 组件事件数据
type ComponentEventData struct {
	ComponentID string `json:"componentId"`
	EventType   string `json:"eventType"`
	Value       string `json:"value"`
}
