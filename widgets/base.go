package widgets

import (
	"fmt"
	"sync/atomic"
	"time"
)

// GlobalUpdateFunc 全局更新函数类型
type GlobalUpdateFunc func(componentID string, html string)

// ISession 会话接口，避免循环依赖
type ISession interface {
	ID() string
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Has(key string) bool
	Clear()
	LastAccessedAt() time.Time
	CreatedAt() time.Time
	AddWidget(widget Widget)
	GetWidgets() []Widget
	ClearWidgets()
	UpdateWidget(componentID string, html string)
}

// Widget 组件接口，所有组件必须实现此接口
type Widget interface {
	// Render 渲染组件为HTML字符串
	Render() string

	// GetID 获取组件唯一标识
	GetID() string

	// GetType 获取组件类型
	GetType() string

	// SetKey 设置组件键值
	SetKey(key string)

	// GetKey 获取组件键值
	GetKey() string

	// OnChange 设置值变更回调函数
	OnChange(callback func(session ISession, event string, value string))

	// IsVisible 检查组件是否可见
	IsVisible() bool
}

// ITriggerCallbacks 触发回调接口
type ITriggerCallbacks interface {
	TriggerCallbacks(session ISession, event string, value string)
}

// BaseWidget 组件基类，提供通用功能
type BaseWidget struct {
	id         string                                               // 唯一标识符
	key        string                                               // 用户定义的键
	widgetType string                                               // 组件类型
	visible    bool                                                 // 可见性标志
	callbacks  []func(session ISession, event string, value string) // 值变更回调函数列表
}

// NewBaseWidget 创建基础组件
func NewBaseWidget(widgetType string) *BaseWidget {
	return &BaseWidget{
		id:         generateID(),
		widgetType: widgetType,
		visible:    true,
		callbacks:  make([]func(session ISession, event string, value string), 0),
	}
}

// GetID 获取组件ID
func (w *BaseWidget) GetID() string {
	return w.id
}

// GetType 获取组件类型
func (w *BaseWidget) GetType() string {
	return w.widgetType
}

// SetKey 设置组件键值
func (w *BaseWidget) SetKey(key string) {
	w.key = key
}

// GetKey 获取组件键值
func (w *BaseWidget) GetKey() string {
	return w.key
}

// OnChange 设置值变更回调函数
func (w *BaseWidget) OnChange(callback func(session ISession, event string, value string)) {
	w.callbacks = append(w.callbacks, callback)
}

// TriggerCallbacks 触发所有回调函数
func (w *BaseWidget) TriggerCallbacks(session ISession, event string, value string) {
	for _, callback := range w.callbacks {
		if callback != nil {
			callback(session, event, value)
		}
	}
}

// UpdateWidget 更新组件并发送局部更新
func (w *BaseWidget) UpdateWidget(session ISession, renderer func() string) string {
	html := renderer()

	// 如果提供了会话上下文，则使用会话的局部更新方法
	if session != nil {
		session.UpdateWidget(w.GetID(), html)
	} else if globalUpdateFunc != nil {
		// 否则使用全局更新函数（向后兼容）
		globalUpdateFunc(w.GetID(), html)
	}

	return html
}

// SetVisible 设置可见性
func (w *BaseWidget) SetVisible(visible bool) {
	w.visible = visible
}

// IsVisible 检查是否可见
func (w *BaseWidget) IsVisible() bool {
	return w.visible
}

// globalUpdateFunc 全局更新函数实例
var globalUpdateFunc GlobalUpdateFunc

// SetGlobalUpdateFunc 设置全局更新函数
func SetGlobalUpdateFunc(updateFunc GlobalUpdateFunc) {
	globalUpdateFunc = updateFunc
}

var widgetIDCounter uint64

// generateID 生成组件唯一ID
func generateID() string {
	id := atomic.AddUint64(&widgetIDCounter, 1)
	return fmt.Sprintf("widget_%d", id)
}
