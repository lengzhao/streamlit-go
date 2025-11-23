package app

import (
	"log"
	"strconv"

	"github.com/lengzhao/streamlit-go/widgets"
)

// Title 添加标题组件
func (a *App) Title(text string, anchor ...string) *widgets.TitleWidget {
	w := widgets.NewTitle(text, anchor...)
	a.AddWidget(w)
	return w
}

// Header 添加Header组件
func (a *App) Header(text string, divider ...bool) *widgets.HeaderWidget {
	w := widgets.NewHeader(text, divider...)
	a.AddWidget(w)
	return w
}

// Subheader 添加Subheader组件
func (a *App) Subheader(text string) *widgets.SubheaderWidget {
	w := widgets.NewSubheader(text)
	a.AddWidget(w)
	return w
}

// Text 添加Text组件
func (a *App) Text(text string) *widgets.TextWidget {
	w := widgets.NewText(text)
	a.AddWidget(w)
	return w
}

// Write 添加Write组件
func (a *App) Write(data interface{}) *widgets.WriteWidget {
	w := widgets.NewWrite(data)
	a.AddWidget(w)
	return w
}

// Button 添加Button组件
// 注意：此方法不直接返回值，需要通过回调函数获取点击事件
func (a *App) Button(label string, key ...string) *widgets.ButtonWidget {
	w := widgets.NewButton(label, key...)
	a.AddWidget(w)
	return w
}

// ButtonWithCallback 添加带回调函数的Button组件
func (a *App) ButtonWithCallback(label string, callback func(session widgets.SessionInterface), key ...string) *widgets.ButtonWidget {
	w := widgets.NewButton(label, key...)
	w.OnChange(func(session widgets.SessionInterface, event string, value string) {
		callback(session)
	})
	a.AddWidget(w)
	return w
}

// TextInput 添加TextInput组件
func (a *App) TextInput(label string, value ...string) *widgets.TextInputWidget {
	var val string
	if len(value) > 0 {
		val = value[0]
	}
	w := widgets.NewTextInput(label, val)
	a.AddWidget(w)
	return w
}

// TextInputWithCallback 添加带回调函数的TextInput组件
func (a *App) TextInputWithCallback(label string, callback func(session widgets.SessionInterface, value string), value ...string) *widgets.TextInputWidget {
	var val string
	if len(value) > 0 {
		val = value[0]
	}
	w := widgets.NewTextInput(label, val)
	w.OnChange(func(session widgets.SessionInterface, event string, value string) {
		callback(session, value)
	})
	a.AddWidget(w)
	return w
}

// NumberInput 添加NumberInput组件
func (a *App) NumberInput(label string, value ...float64) *widgets.NumberInputWidget {
	var val float64
	if len(value) > 0 {
		val = value[0]
	}
	w := widgets.NewNumberInput(label, val)
	a.AddWidget(w)
	return w
}

// NumberInputWithCallback 添加带回调函数的NumberInput组件
func (a *App) NumberInputWithCallback(label string, callback func(session widgets.SessionInterface, value float64), value ...float64) *widgets.NumberInputWidget {
	var val float64
	if len(value) > 0 {
		val = value[0]
	}
	w := widgets.NewNumberInput(label, val)
	w.OnChange(func(session widgets.SessionInterface, event string, value string) {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			callback(session, f)
		} else {
			log.Printf("Failed to parse number input value: %s", value)
		}
	})
	a.AddWidget(w)
	return w
}

// Container 添加Container组件
func (a *App) Container(border ...bool) *widgets.ContainerWidget {
	var b bool
	if len(border) > 0 {
		b = border[0]
	}
	w := widgets.NewContainer(b)
	a.AddWidget(w)
	return w
}

// Columns 添加Columns组件
func (a *App) Columns(ratios ...int) []*widgets.Column {
	w := widgets.NewColumns(ratios...)
	a.AddWidget(w)
	return w.GetColumns()
}

// Sidebar 添加Sidebar组件
func (a *App) Sidebar(expanded ...bool) *widgets.SidebarWidget {
	var e bool
	if len(expanded) > 0 {
		e = expanded[0]
	}
	w := widgets.NewSidebar(e)
	a.AddWidget(w)
	return w
}

// Expander 添加Expander组件
func (a *App) Expander(label string, expanded ...bool) *widgets.ExpanderWidget {
	var e bool
	if len(expanded) > 0 {
		e = expanded[0]
	}
	w := widgets.NewExpander(label, e)
	a.AddWidget(w)
	return w
}

// Table 添加Table组件
func (a *App) Table(data interface{}) *widgets.TableWidget {
	w := widgets.NewTable(data)
	a.AddWidget(w)
	return w
}

// DataFrame 添加DataFrame组件
func (a *App) DataFrame(data interface{}) *widgets.DataFrameWidget {
	w := widgets.NewDataFrame(data)
	a.AddWidget(w)
	return w
}

// Metric 添加Metric组件
func (a *App) Metric(label string, value interface{}) *widgets.MetricWidget {
	w := widgets.NewMetric(label, value)
	a.AddWidget(w)
	return w
}
