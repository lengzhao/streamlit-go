package state

import (
	"sync"
	"time"

	"github.com/lengzhao/streamlit-go/widgets"
)

// Session 会话结构，存储单个用户的会话状态
type Session struct {
	id             string           // 会话唯一标识
	widgets        []widgets.Widget // 会话私有组件
	widgetsMutex   sync.RWMutex     // 组件队列锁
	createdAt      time.Time        // 创建时间
	lastAccessedAt time.Time        // 最后访问时间
	mutex          sync.RWMutex     // 读写锁，保护并发访问
}

// NewSession 创建新的会话
func NewSession(id string) *Session {
	now := time.Now()
	return &Session{
		id:             id,
		widgets:        make([]widgets.Widget, 0),
		widgetsMutex:   sync.RWMutex{},
		createdAt:      now,
		lastAccessedAt: now,
		mutex:          sync.RWMutex{},
	}
}

// ID 返回会话ID
func (s *Session) ID() string {
	return s.id
}

// LastAccessedAt 返回最后访问时间
func (s *Session) LastAccessedAt() time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.lastAccessedAt
}

// CreatedAt 返回创建时间
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

// AddWidget 添加组件到会话
func (s *Session) AddWidget(widget widgets.Widget) {
	s.widgetsMutex.Lock()
	defer s.widgetsMutex.Unlock()

	s.widgets = append(s.widgets, widget)
}

func (s *Session) SetWidget(widget widgets.Widget) {
	s.widgetsMutex.Lock()
	defer s.widgetsMutex.Unlock()
	for i, w := range s.widgets {
		if w.GetID() == widget.GetID() {
			s.widgets[i] = widget
			return
		}
	}
}

// GetWidgets 获取会话组件
func (s *Session) GetWidgets() []widgets.Widget {
	s.widgetsMutex.RLock()
	defer s.widgetsMutex.RUnlock()

	// 返回副本，避免并发问题
	widgetsCopy := make([]widgets.Widget, len(s.widgets))
	copy(widgetsCopy, s.widgets)
	return widgetsCopy
}

// ClearWidgets 清空会话组件
func (s *Session) ClearWidgets() {
	s.widgetsMutex.Lock()
	defer s.widgetsMutex.Unlock()

	s.widgets = make([]widgets.Widget, 0)
}

// RemoveWidget 从会话中移除指定ID的组件
func (s *Session) RemoveWidget(componentID string) {
	s.widgetsMutex.Lock()
	defer s.widgetsMutex.Unlock()

	// 查找并移除指定ID的组件
	for i, widget := range s.widgets {
		if widget.GetID() == componentID {
			// 从切片中移除该组件
			s.widgets = append(s.widgets[:i], s.widgets[i+1:]...)
			break
		}
	}
}

// DeleteWidget 删除组件（占位方法，实际实现在前端）
func (s *Session) DeleteWidget(componentID string) {
	// 从会话中移除组件
	s.RemoveWidget(componentID)

	// 占位方法，实际删除逻辑由前端处理
}

// LastAccessedAtStr 返回最后访问时间的字符串表示
func (s *Session) LastAccessedAtStr() string {
	return s.lastAccessedAt.Format(time.RFC3339)
}

// CreatedAtStr 返回创建时间的字符串表示
func (s *Session) CreatedAtStr() string {
	return s.createdAt.Format(time.RFC3339)
}
