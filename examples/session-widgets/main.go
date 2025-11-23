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
		app.WithPort(8503),
	)

	// 添加全局组件（所有用户共享）
	st.Title("🌐 全局标题 - 所有用户共享")
	st.Header("这是全局内容", true)
	st.Text("所有用户都会看到这段文字")

	// 为不同会话添加私有组件
	user1SessionID := "user-1-session"
	user2SessionID := "user-2-session"

	// 为用户1添加私有组件
	user1Session := st.GetStateManager().GetSession(user1SessionID)
	user1Title := widgets.NewTitle("👤 用户1的私有标题")
	user1Session.AddWidget(user1Title)
	user1Text := widgets.NewText("这是用户1的私有内容")
	user1Session.AddWidget(user1Text)
	user1Counter := widgets.NewWrite("用户1计数器: 0")
	user1Session.AddWidget(user1Counter)

	// 为用户2添加私有组件
	user2Session := st.GetStateManager().GetSession(user2SessionID)
	user2Title := widgets.NewTitle("👤 用户2的私有标题")
	user2Session.AddWidget(user2Title)
	user2Text := widgets.NewText("这是用户2的私有内容")
	user2Session.AddWidget(user2Text)
	user2Counter := widgets.NewWrite("用户2计数器: 0")
	user2Session.AddWidget(user2Counter)

	// 为用户1创建按钮
	user1Button := widgets.NewButton("用户1按钮")
	user1Session.AddWidget(user1Button)
	user1Count := 0
	user1Button.OnChange(func(session widgets.SessionInterface, event string, value string) {
		user1Count++
		user1Counter.SetData("用户1计数器: " + fmt.Sprintf("%d", user1Count))
	})

	// 为用户2创建按钮
	user2Button := widgets.NewButton("用户2按钮")
	user2Session.AddWidget(user2Button)
	user2Count := 0
	user2Button.OnChange(func(session widgets.SessionInterface, event string, value string) {
		user2Count++
		user2Counter.SetData("用户2计数器: " + fmt.Sprintf("%d", user2Count))
	})

	// 添加说明文本
	st.Subheader("👥 会话隔离演示")
	st.Text("不同用户只能看到自己的私有组件和全局组件")
	st.Text("请使用不同的浏览器标签页或设备，分别访问以下URL进行测试：")
	st.Text("用户1: http://localhost:8503?sessionId=user-1-session")
	st.Text("用户2: http://localhost:8503?sessionId=user-2-session")
	st.Text("默认: http://localhost:8503?sessionId=default-session-id")

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
