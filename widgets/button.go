package widgets

import (
	"fmt"
)

// ButtonWidget 按钮组件
type ButtonWidget struct {
	*BaseWidget
	label string
}

// NewButton 创建新的按钮组件
func NewButton(label string) *ButtonWidget {
	w := &ButtonWidget{
		BaseWidget: NewBaseWidget("button"),
		label:      label,
	}
	// 移除SetKey调用，因为key参数将被忽略
	return w
}

// Render 渲染按钮组件为HTML
func (w *ButtonWidget) Render() string {
	id := w.GetID()
	return fmt.Sprintf("<button class=\"st-button\" data-widget-id=\"%s\" id=\"%s\" data-event-type=\"click\">%s</button>", id, id, w.label)
}

// SetValue 设置按钮值（触发点击事件）
func (w *ButtonWidget) SetValue(session ISession) {
	w.TriggerCallbacks(session, "click", "")
}
