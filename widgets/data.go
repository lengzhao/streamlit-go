package widgets

import (
	"fmt"
	"html"
	"reflect"
)

// TableWidget 表格组件
type TableWidget struct {
	*BaseWidget
	data interface{}
}

// NewTable 创建新的表格组件
func NewTable(data interface{}) *TableWidget {
	w := &TableWidget{
		BaseWidget: NewBaseWidget("table"),
		data:       data,
	}

	return w
}

// Render 渲染表格组件为HTML
func (w *TableWidget) Render() string {
	// 简单实现，支持字符串切片
	switch v := w.data.(type) {
	case []string:
		rows := ""
		for _, item := range v {
			rows += fmt.Sprintf("<tr><td>%s</td></tr>", html.EscapeString(item))
		}
		return fmt.Sprintf("<table class=\"st-table\" data-widget-id=\"%s\"><tbody>%s</tbody></table>", w.GetID(), rows)
	default:
		return fmt.Sprintf("<div class=\"st-table\" data-widget-id=\"%s\">%v</div>", w.GetID(), w.data)
	}
}

// DataFrameWidget 数据框组件
type DataFrameWidget struct {
	*BaseWidget
	data interface{}
}

// NewDataFrame 创建新的数据框组件
func NewDataFrame(data interface{}) *DataFrameWidget {
	w := &DataFrameWidget{
		BaseWidget: NewBaseWidget("dataframe"),
		data:       data,
	}

	return w
}

// Render 渲染数据框组件为HTML
func (w *DataFrameWidget) Render() string {
	// 简单实现，支持map[string]interface{}
	switch v := w.data.(type) {
	case map[string]interface{}:
		rows := ""
		for key, value := range v {
			rows += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", html.EscapeString(key), value)
		}
		return fmt.Sprintf("<table class=\"st-dataframe\" data-widget-id=\"%s\"><tbody>%s</tbody></table>", w.GetID(), rows)
	case map[string]string:
		rows := ""
		for key, value := range v {
			rows += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", html.EscapeString(key), html.EscapeString(value))
		}
		return fmt.Sprintf("<table class=\"st-dataframe\" data-widget-id=\"%s\"><tbody>%s</tbody></table>", w.GetID(), rows)
	default:
		// 使用反射来处理其他类型
		val := reflect.ValueOf(w.data)
		if val.Kind() == reflect.Struct {
			rows := ""
			t := val.Type()
			for i := 0; i < val.NumField(); i++ {
				field := t.Field(i)
				value := val.Field(i)
				rows += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", html.EscapeString(field.Name), value.Interface())
			}
			return fmt.Sprintf("<table class=\"st-dataframe\" data-widget-id=\"%s\"><tbody>%s</tbody></table>", w.GetID(), rows)
		}
		return fmt.Sprintf("<div class=\"st-dataframe\" data-widget-id=\"%s\">%v</div>", w.GetID(), w.data)
	}
}

// MetricWidget 指标组件
type MetricWidget struct {
	*BaseWidget
	label string
	value interface{}
	delta string
}

// NewMetric 创建新的指标组件
func NewMetric(label string, value interface{}) *MetricWidget {
	w := &MetricWidget{
		BaseWidget: NewBaseWidget("metric"),
		label:      label,
		value:      value,
	}
	return w
}

// SetDelta 设置指标变化值
func (w *MetricWidget) SetDelta(delta string) {
	w.delta = delta
}

// Render 渲染指标组件为HTML
func (w *MetricWidget) Render() string {
	deltaHTML := ""
	if w.delta != "" {
		deltaHTML = fmt.Sprintf("<div class=\"st-metric-delta\">%s</div>", html.EscapeString(w.delta))
	}
	return fmt.Sprintf("<div class=\"st-metric\" data-widget-id=\"%s\"><div class=\"st-metric-label\">%s</div><div class=\"st-metric-value\">%v</div>%s</div>",
		w.GetID(), html.EscapeString(w.label), w.value, deltaHTML)
}
