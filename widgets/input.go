package widgets

import (
	"fmt"
	"html"
	"strconv"
)

// TextInputWidget 文本输入组件
type TextInputWidget struct {
	*BaseWidget
	label       string
	value       string
	placeholder string
}

// NewTextInput 创建新的文本输入组件
func NewTextInput(label string, value string, key ...string) *TextInputWidget {
	w := &TextInputWidget{
		BaseWidget: NewBaseWidget("text_input"),
		label:      label,
		value:      value,
	}
	// 移除SetKey调用，因为key参数将被忽略
	return w
}

// Render 渲染文本输入组件为HTML
func (w *TextInputWidget) Render() string {
	placeholderAttr := ""
	if w.placeholder != "" {
		placeholderAttr = fmt.Sprintf(" placeholder=\"%s\"", html.EscapeString(w.placeholder))
	}
	return fmt.Sprintf("<div class=\"st-text-input-container\" data-widget-id=\"%s\"><label>%s</label><input type=\"text\" class=\"st-text-input\" data-widget-id=\"%s\" data-event-type=\"input\" value=\"%s\"%s></div>",
		w.GetID(), html.EscapeString(w.label), w.GetID(), html.EscapeString(w.value), placeholderAttr)
}

// SetValue 设置文本输入值
func (w *TextInputWidget) SetValue(session ISession, value string) {
	w.value = value
	w.TriggerCallbacks(session, "input", value)
}

// GetValue 获取文本输入值
func (w *TextInputWidget) GetValue() string {
	return w.value
}

// SetPlaceholder 设置占位符
func (w *TextInputWidget) SetPlaceholder(placeholder string) {
	w.placeholder = placeholder
}

// NumberInputWidget 数字输入组件
type NumberInputWidget struct {
	*BaseWidget
	label string
	value float64
	step  float64
}

// NewNumberInput 创建新的数字输入组件
func NewNumberInput(label string, value float64) *NumberInputWidget {
	w := &NumberInputWidget{
		BaseWidget: NewBaseWidget("number_input"),
		label:      label,
		value:      value,
		step:       1,
	}
	// 移除SetKey调用，因为key参数将被忽略
	return w
}

// Render 渲染数字输入组件为HTML
func (w *NumberInputWidget) Render() string {
	return fmt.Sprintf("<div class=\"st-number-input-container\" data-widget-id=\"%s\"><label>%s</label><input type=\"number\" class=\"st-number-input\" data-widget-id=\"%s\" data-event-type=\"input\" value=\"%g\" step=\"%g\"></div>",
		w.GetID(), html.EscapeString(w.label), w.GetID(), w.value, w.step)
}

// SetValue 设置数字输入值
func (w *NumberInputWidget) SetValue(session ISession, value float64) {
	w.value = value
	w.TriggerCallbacks(session, "input", strconv.FormatFloat(value, 'g', -1, 64))
}

// GetValue 获取数字输入值
func (w *NumberInputWidget) GetValue() float64 {
	return w.value
}

// SetStep 设置步长
func (w *NumberInputWidget) SetStep(step float64) {
	w.step = step
}
