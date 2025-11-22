package ui

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/lengzhao/streamlit-go/ptemplate"
	"github.com/lengzhao/streamlit-go/widgets"
)

// Renderer UI渲染器
type Renderer struct {
}

// NewRenderer 创建新的UI渲染器
func NewRenderer() *Renderer {
	return &Renderer{}
}

// RenderWidgets 渲染所有组件为HTML
func (r *Renderer) RenderWidgets(widgets []widgets.Widget) string {
	html := ""
	for _, widget := range widgets {
		if widget.IsVisible() {
			html += widget.Render()
		}
	}
	return html
}

// RenderPage 渲染完整页面
func (r *Renderer) RenderPage(title string, widgets []widgets.Widget) (string, error) {
	widgetsHTML := r.RenderWidgets(widgets)

	// 从模板包中获取页面模板
	tmpl, err := ptemplate.GetPageTemplate()
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := map[string]interface{}{
		"Title":   title,
		"Content": template.HTML(widgetsHTML),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
