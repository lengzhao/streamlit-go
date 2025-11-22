package widgets

import (
"fmt"
)

// ContainerWidget 容器组件
type ContainerWidget struct {
	*BaseWidget
	border bool
	children []Widget
}

// NewContainer 创建新的容器组件
func NewContainer(border bool, key ...string) *ContainerWidget {
	w := &ContainerWidget{
		BaseWidget: NewBaseWidget("container"),
		border:     border,
		children:   make([]Widget, 0),
	}
	if len(key) > 0 {
		w.SetKey(key[0])
	}
	return w
}

// AddChild 添加子组件
func (w *ContainerWidget) AddChild(child Widget) {
	w.children = append(w.children, child)
}

// Render 渲染容器组件为HTML
func (w *ContainerWidget) Render() string {
	borderClass := ""
	if w.border {
		borderClass = " st-container-with-border"
	}
	
	childrenHTML := ""
	for _, child := range w.children {
		childrenHTML += child.Render()
	}
	
	return fmt.Sprintf("<div class=\"st-container%s\" data-widget-id=\"%s\">%s</div>", borderClass, w.GetID(), childrenHTML)
}

// Column 列组件
type Column struct {
	*BaseWidget
	ratio    int
	children []Widget
}

// NewColumn 创建新的列组件
func NewColumn(ratio int) *Column {
	return &Column{
		BaseWidget: NewBaseWidget("column"),
		ratio:      ratio,
		children:   make([]Widget, 0),
	}
}

// AddChild 添加子组件
func (c *Column) AddChild(child Widget) {
	c.children = append(c.children, child)
}

// Render 渲染列组件为HTML
func (c *Column) Render() string {
	childrenHTML := ""
	for _, child := range c.children {
		childrenHTML += child.Render()
	}
	
	return fmt.Sprintf("<div class=\"st-column\" style=\"flex: %d\" data-widget-id=\"%s\">%s</div>", c.ratio, c.GetID(), childrenHTML)
}

// ColumnsWidget 列布局组件
type ColumnsWidget struct {
	*BaseWidget
	columns []*Column
	ratios  []int
}

// NewColumns 创建新的列布局组件
func NewColumns(ratios ...int) *ColumnsWidget {
	// 如果没有提供比率，默认为相等比率
	if len(ratios) == 0 {
		ratios = []int{1, 1}
	}
	
	columns := make([]*Column, len(ratios))
	for i, ratio := range ratios {
		columns[i] = NewColumn(ratio)
	}
	
	return &ColumnsWidget{
		BaseWidget: NewBaseWidget("columns"),
		columns:    columns,
		ratios:     ratios,
	}
}

// GetColumns 获取列组件数组
func (w *ColumnsWidget) GetColumns() []*Column {
	return w.columns
}

// Render 渲染列布局组件为HTML
func (w *ColumnsWidget) Render() string {
	columnsHTML := ""
	for _, column := range w.columns {
		columnsHTML += column.Render()
	}
	
	return fmt.Sprintf("<div class=\"st-columns\" data-widget-id=\"%s\">%s</div>", w.GetID(), columnsHTML)
}

// SidebarWidget 侧边栏组件
type SidebarWidget struct {
	*BaseWidget
	expanded bool
	children []Widget
}

// NewSidebar 创建新的侧边栏组件
func NewSidebar(expanded bool, key ...string) *SidebarWidget {
	w := &SidebarWidget{
		BaseWidget: NewBaseWidget("sidebar"),
		expanded:   expanded,
		children:   make([]Widget, 0),
	}
	if len(key) > 0 {
		w.SetKey(key[0])
	}
	return w
}

// AddChild 添加子组件
func (w *SidebarWidget) AddChild(child Widget) {
	w.children = append(w.children, child)
}

// Render 渲染侧边栏组件为HTML
func (w *SidebarWidget) Render() string {
	expandedClass := ""
	if w.expanded {
		expandedClass = " st-sidebar-expanded"
	}
	
	childrenHTML := ""
	for _, child := range w.children {
		childrenHTML += child.Render()
	}
	
	return fmt.Sprintf("<div class=\"st-sidebar%s\" data-widget-id=\"%s\">%s</div>", expandedClass, w.GetID(), childrenHTML)
}

// ExpanderWidget 可展开组件
type ExpanderWidget struct {
	*BaseWidget
	label    string
	expanded bool
	children []Widget
}

// NewExpander 创建新的可展开组件
func NewExpander(label string, expanded bool, key ...string) *ExpanderWidget {
	w := &ExpanderWidget{
		BaseWidget: NewBaseWidget("expander"),
		label:      label,
		expanded:   expanded,
		children:   make([]Widget, 0),
	}
	if len(key) > 0 {
		w.SetKey(key[0])
	}
	return w
}

// AddChild 添加子组件
func (w *ExpanderWidget) AddChild(child Widget) {
	w.children = append(w.children, child)
}

// Render 渲染可展开组件为HTML
func (w *ExpanderWidget) Render() string {
	expandedClass := ""
	if w.expanded {
		expandedClass = " st-expander-expanded"
	}
	
	childrenHTML := ""
	for _, child := range w.children {
		childrenHTML += child.Render()
	}
	
	return fmt.Sprintf("<div class=\"st-expander%s\" data-widget-id=\"%s\"><div class=\"st-expander-header\">%s</div><div class=\"st-expander-content\">%s</div></div>", 
expandedClass, w.GetID(), w.label, childrenHTML)
}
