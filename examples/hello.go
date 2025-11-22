package main

import (
"fmt"
"log"
"os"
"os/signal"
"syscall"
"time"

"github.com/lengzhao/streamlit-go/app"
)

func main() {
	// 创建应用实例
	st := app.New(
app.WithTitle("我的第一个Streamlit Go应用"),
app.WithPort(8501),
)

	// 添加各种组件进行测试
	st.Title("🚀 欢迎使用Streamlit Go")
	st.Header("这是一个全功能示例", true)

	// 文本组件
	st.Subheader("📝 文本组件")
	st.Text("这是普通文本")
	st.Write("这是Write组件，可以展示各种数据类型")
	st.Write(42)
	st.Write(3.14159)
	st.Write(true)

	// 指标组件
	st.Subheader("�� 指标展示")
	metric1 := st.Metric("总用户数", 1234)
	metric1.SetDelta("+12%")

	metric2 := st.Metric("活跃用户", 567)
	metric2.SetDelta("+5%")

	metric3 := st.Metric("收入", "$89,432")
	metric3.SetDelta("-2.3%")

	// 输入组件
	st.Subheader("⌨️ 输入组件")
	// 先创建一个WriteWidget来显示输入的姓名
	nameOutput := st.Write("")
	nameInput := st.TextInputWithCallback("姓名", func(name string) {
if name != "" {
nameOutput.SetData("您好，" + name + "！")
}
}, "")
	nameInput.SetPlaceholder("请输入您的姓名")

	// 先创建一个WriteWidget来显示输入的年龄
	ageOutput := st.Write("")
	st.NumberInputWithCallback("您的年龄", func(age float64) {
log.Printf("年龄 changed! Value: %v", age)
ageOutput.SetData(fmt.Sprintf("您的年龄是: %.0f", age))
}, 25)

	// 按钮
	// 先创建一个WriteWidget来显示按钮点击消息
	buttonOutput := st.Write("")
	st.ButtonWithCallback("👍 点击我", func() {
		buttonOutput.SetData(fmt.Sprintf("按钮被点击了！%v", time.Now().Format("2006-01-02 15:04:05")))
	})

	// 数据展示
	st.Subheader("📊 数据展示")

	// 简单表格
	data := []string{"苹果", "香蕉", "橙子"}
	st.Table(data)

	// Map数据
	mapData := map[string]interface{}{
		"名称": "Streamlit Go",
		"版本": "0.1.0",
		"语言": "Golang",
	}
	st.DataFrame(mapData)

	// 布局组件
	st.Subheader("📏 布局组件")

	container := st.Container(true)
	containerText := st.Text("这是一个带边框的容器")
	container.AddChild(containerText)

	expander := st.Expander("🔍 点击展开查看更多", false)
	expanderText := st.Text("这是隐藏的内容，点击标题可以展开或折叠")
	expander.AddChild(expanderText)

	log.Println("应用创建成功")
	log.Println("请在浏览器中访问 http://localhost:8501")

	// 设置信号处理，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 在单独的goroutine中启动应用
	go func() {
		if err := st.Run(); err != nil {
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
