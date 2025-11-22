package main

import (
"fmt"
"log"
"os"
"os/signal"
"syscall"

"github.com/lengzhao/streamlit-go/app"
)

func main() {
	// 创建应用实例
	st := app.New(
app.WithTitle("Streamlit Go"),
app.WithPort(8503),
)

	// 添加组件
	st.Title("欢迎使用 Streamlit Go")
	st.Header("这是一个基础示例", true)
	st.Subheader("文本组件")
	st.Text("这是普通文本")
	st.Write("这是Write组件")
	st.Write(42)
	st.Write(3.14)
	st.Write(true)

	// 按钮示例
	buttonOutput := st.Write("")
	st.ButtonWithCallback("点击我", func() {
		buttonOutput.SetData("按钮被点击了！")
	})

	// 输入示例
	inputOutput := st.Write("")
	textInput := st.TextInputWithCallback("输入文本:", func(value string) {
inputOutput.SetData("您输入的是: " + value)
}, "")
	textInput.SetPlaceholder("请输入文本")

	// 打印组件信息
	widgets := st.GetWidgets()
	fmt.Printf("创建了 %d 个组件:\n", len(widgets))
	for i, w := range widgets {
		fmt.Printf("%d. 类型: %s, ID: %s\n", i+1, w.GetType(), w.GetID())
	}

	// 测试会话状态
	session := st.Session()
	session.Set("test_key", "test_value")
	if val, ok := session.Get("test_key"); ok {
		fmt.Printf("\n会话状态测试成功: %v\n", val)
	}

	fmt.Println("\n基础框架测试完成！")

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
	fmt.Println("\n收到中断信号，关闭中...")

	// 优雅关闭
	if err := st.Stop(); err != nil {
		log.Printf("关闭时错误: %v", err)
	}

	fmt.Println("应用已成功停止")
}
