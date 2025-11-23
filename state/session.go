package state

import (
	"sync"
	"time"

	"github.com/lengzhao/streamlit-go/widgets"
)

// Session 会话结构，存储单个用户的会话状态
type Session struct {
	id             string                 // 会话唯一标识
	state          map[string]interface{} // 状态存储
	widgets        []widgets.Widget       // 会话私有组件
	widgetsMutex   sync.RWMutex           // 组件队列锁
	createdAt      time.Time              // 创建时间
	lastAccessedAt time.Time              // 最后访问时间
	mutex          sync.RWMutex           // 读写锁，保护并发访问
}

// NewSession 创建新的会话
func NewSession(id string) *Session {
	now := time.Now()
	return &Session{
		id:             id,
		state:          make(map[string]interface{}),
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

// Get 获取状态值
func (s *Session) Get(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.lastAccessedAt = time.Now()
	value, exists := s.state[key]
	return value, exists
}

// Set 设置状态值
func (s *Session) Set(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.state[key] = value
	s.lastAccessedAt = time.Now()
}

// Delete 删除状态值
func (s *Session) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.state, key)
	s.lastAccessedAt = time.Now()
}

// Has 检查键是否存在
func (s *Session) Has(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.state[key]
	return exists
}

// Clear 清空所有状态
func (s *Session) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.state = make(map[string]interface{})
	s.lastAccessedAt = time.Now()
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

// LastAccessedAtStr 返回最后访问时间的字符串表示
func (s *Session) LastAccessedAtStr() string {
	return s.lastAccessedAt.Format(time.RFC3339)
}

// CreatedAtStr 返回创建时间的字符串表示
func (s *Session) CreatedAtStr() string {
	return s.createdAt.Format(time.RFC3339)
}
