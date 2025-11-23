package widgets

import (
	"fmt"
	"html"
)

// TitleWidget 标题组件
type TitleWidget struct {
	*BaseWidget
	text   string
	anchor string
}

// NewTitle 创建新的标题组件
func NewTitle(text string, anchor ...string) *TitleWidget {
	w := &TitleWidget{
		BaseWidget: NewBaseWidget("title"),
		text:       text,
	}
	if len(anchor) > 0 {
		w.anchor = anchor[0]
	}
	return w
}

// SetText 设置标题文本
func (w *TitleWidget) SetText(text string) {
	w.text = text
	// 触发更新（注意：这里没有会话上下文，可能需要在其他地方处理更新）
	w.UpdateWidget(nil, func() string {
		return w.Render()
	})
}

// Render 渲染标题组件为HTML
func (w *TitleWidget) Render() string {
	anchorAttr := ""
	if w.anchor != "" {
		anchorAttr = fmt.Sprintf(" id=\"%s\"", html.EscapeString(w.anchor))
	}
	return fmt.Sprintf("<h1 class=\"st-title\" data-widget-id=\"%s\"%s>%s</h1>", w.GetID(), anchorAttr, html.EscapeString(w.text))
}

// HeaderWidget 二级标题组件
type HeaderWidget struct {
	*BaseWidget
	text    string
	divider bool
}

// NewHeader 创建新的二级标题组件
func NewHeader(text string, divider ...bool) *HeaderWidget {
	w := &HeaderWidget{
		BaseWidget: NewBaseWidget("header"),
		text:       text,
	}
	if len(divider) > 0 {
		w.divider = divider[0]
	}
	return w
}

// Render 渲染二级标题组件为HTML
func (w *HeaderWidget) Render() string {
	dividerClass := ""
	if w.divider {
		dividerClass = " st-header-with-divider"
	}
	return fmt.Sprintf("<h2 class=\"st-header%s\" data-widget-id=\"%s\">%s</h2>", dividerClass, w.GetID(), html.EscapeString(w.text))
}

// SubheaderWidget 三级标题组件
type SubheaderWidget struct {
	*BaseWidget
	text string
}

// NewSubheader 创建新的三级标题组件
func NewSubheader(text string) *SubheaderWidget {
	return &SubheaderWidget{
		BaseWidget: NewBaseWidget("subheader"),
		text:       text,
	}
}

// Render 渲染三级标题组件为HTML
func (w *SubheaderWidget) Render() string {
	return fmt.Sprintf("<h3 class=\"st-subheader\" data-widget-id=\"%s\">%s</h3>", w.GetID(), html.EscapeString(w.text))
}

// TextWidget 文本组件
type TextWidget struct {
	*BaseWidget
	text string
}

// NewText 创建新的文本组件
func NewText(text string) *TextWidget {
	return &TextWidget{
		BaseWidget: NewBaseWidget("text"),
		text:       text,
	}
}

// Render 渲染文本组件为HTML
func (w *TextWidget) Render() string {
	return fmt.Sprintf("<div class=\"st-text\" data-widget-id=\"%s\">%s</div>", w.GetID(), html.EscapeString(w.text))
}

// WriteWidget 通用数据展示组件
type WriteWidget struct {
	*BaseWidget
	data interface{}
}

// NewWrite 创建新的通用数据展示组件
func NewWrite(data interface{}) *WriteWidget {
	return &WriteWidget{
		BaseWidget: NewBaseWidget("write"),
		data:       data,
	}
}

// SetData 设置数据
func (w *WriteWidget) SetData(data interface{}) {
	w.data = data
	// 触发更新（注意：这里没有会话上下文，可能需要在其他地方处理更新）
	w.UpdateWidget(nil, func() string {
		return w.Render()
	})
}

// GetData 获取数据
func (w *WriteWidget) GetData() interface{} {
	return w.data
}

// Render 渲染通用数据展示组件为HTML
func (w *WriteWidget) Render() string {
	return fmt.Sprintf("<div class=\"st-write\" data-widget-id=\"%s\">%v</div>", w.GetID(), w.data)
}
