package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lengzhao/streamlit-go/core"
	"github.com/lengzhao/streamlit-go/widgets"
)

func main() {
	// 创建应用实例
	st := core.NewService(
		core.WithTitle("我的第一个Streamlit Go应用"),
		core.WithPort(8504),
	)

	// 添加各种组件进行测试
	st.Title("🚀 欢迎使用Streamlit Go")
	st.Header("这是一个全功能示例", true)

	st.AddWidget(widgets.NewSubheader("📝 文本组件"))
	st.Text("这是普通文本")
	st.AddWidget(widgets.NewWrite("这是Write组件，可以展示各种数据类型"))
	st.AddWidget(widgets.NewWrite(42))
	st.AddWidget(widgets.NewWrite(42))
	st.AddWidget(widgets.NewWrite(3.14159))
	st.AddWidget(widgets.NewWrite(true))

	// 指标组件
	st.AddWidget(widgets.NewSubheader("📊 指标展示"))
	metric1 := widgets.NewMetric("总用户数", 1234)
	metric1.SetDelta("+12%")
	st.AddWidget(metric1)

	metric2 := widgets.NewMetric("活跃用户", 567)
	metric2.SetDelta("+5%")
	st.AddWidget(metric2)

	metric3 := widgets.NewMetric("收入", "$89,432")
	metric3.SetDelta("-2.3%")
	st.AddWidget(metric3)

	// 数据展示
	st.AddWidget(widgets.NewSubheader("📈 数据展示"))

	// 简单表格
	data := []string{"苹果", "香蕉", "橙子"}
	st.AddWidget(widgets.NewTable(data))

	// Map数据
	mapData := map[string]interface{}{
		"名称": "Streamlit Go",
		"版本": "0.1.0",
		"语言": "Golang",
	}
	st.AddWidget(widgets.NewDataFrame(mapData))

	// 布局组件
	st.AddWidget(widgets.NewSubheader("📏 布局组件"))

	container := widgets.NewContainer(true)
	containerText := widgets.NewText("这是一个带边框的容器")
	container.AddChild(containerText)
	st.AddWidget(container)

	expander := widgets.NewExpander("🔍 点击展开查看更多", false)
	expanderText := widgets.NewText("这是隐藏的内容，点击标题可以展开或折叠")
	expander.AddChild(expanderText)
	st.AddWidget(expander)

	// 会话特定Widgets示例
	st.AddWidget(widgets.NewSubheader("👥 会话特定Widgets示例"))
	st.Text("以下组件演示了如何为不同用户创建独立的Widgets")

	log.Println("应用创建成功")
	log.Println("请在浏览器中访问 http://localhost:8504")

	// 设置信号处理，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 在单独的goroutine中启动应用
	go func() {
		if err := st.Start(); err != nil {
			log.Printf("服务器错误: %v", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("\n收到中断信号，关闭中...")

	// 优雅关闭
	if err := st.Stop(); err != nil {
		log.Printf("关闭时错误: %v", err)
	}

	log.Println("应用已成功停止")
}
