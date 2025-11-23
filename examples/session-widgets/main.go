package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lengzhao/streamlit-go/app"
	"github.com/lengzhao/streamlit-go/widgets"
)

func main() {
	// 创建应用实例
	st := app.New(
		app.WithTitle("会话Widgets示例"),
		app.WithPort(8504), // 修改端口号避免冲突
	)

	// 添加全局组件（所有用户共享）
	st.Title("🌐 全局标题 - 所有用户共享")
	st.Header("这是全局内容", true)
	st.Text("所有用户都会看到这段文字")

	// 为用户1创建按钮
	user1Button := widgets.NewButton("用户1按钮")

	user1Count := 0
	user1Button.OnChange(func(session widgets.SessionInterface, event string, value string) {
		log.Println("Button clicked by user:", session.ID())
		user1Count++
		stat := widgets.NewText("用户计数器: " + fmt.Sprintf("%d", user1Count))
		session.AddWidget(stat)
	})
	st.AddWidget(user1Button)

	log.Println("应用创建成功")
	log.Println("请在浏览器中访问 http://localhost:8503 查看应用")

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
