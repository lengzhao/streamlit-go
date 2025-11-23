package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lengzhao/streamlit-go/app"
	"github.com/lengzhao/streamlit-go/state"
	"github.com/lengzhao/streamlit-go/widgets"
)

func main() {
	// 创建应用实例
	st := app.New(
		app.WithTitle("登录演示"),
		app.WithPort(8507), // 更改端口为8507
	)

	// 添加全局登录组件
	st.Title("🔐 用户登录")
	loginInput := st.TextInput("用户名", "")
	loginButton := st.Button("登录")

	// 设置登录回调
	st.SetLoginCallback(func(session *state.Session) {
		log.Printf("用户 %s 登录成功", session.ID())

		// 清空之前的组件
		session.ClearWidgets()

		// 添加欢迎信息
		title := widgets.NewTitle("欢迎，" + session.ID())
		session.AddWidget(title)

		text := widgets.NewText("这是您的个人仪表板")
		session.AddWidget(text)

		// 添加计数器
		counter := widgets.NewWrite("计数器: 0")
		session.AddWidget(counter)

		// 添加按钮
		button := widgets.NewButton("增加计数")
		button.OnChange(func(session widgets.SessionInterface, event string, value string) {
			// 获取当前计数
			currentData := counter.GetData()
			currentText := fmt.Sprintf("%v", currentData)
			// 简单解析计数（实际应用中应更健壮）
			count := 0
			if len(currentText) > 8 {
				// "计数器: 0" 中的数字部分
				count = int(currentText[6] - '0')
			}
			count++
			counter.SetData("计数器: " + fmt.Sprintf("%d", count))
		})
		session.AddWidget(button)

		// 触发UI更新
		// 通过更新组件来触发UI更新，而不是直接向session发送消息
		title.SetText("欢迎，" + session.ID())
	})

	// 设置登录按钮回调
	loginButton.OnChange(func(session widgets.SessionInterface, event string, value string) {
		log.Printf("用户 %s 尝试登录", loginInput.GetValue())
		// 注意：这里我们需要将 SessionInterface 转换为 *state.Session
		stateSession, ok := session.(*state.Session)
		if !ok {
			log.Printf("无法转换会话类型")
			return
		}

		username := loginInput.GetValue()
		if username != "" {
			log.Printf("用户 %s 尝试登录", username)
			// 这里可以添加实际的认证逻辑
			// 模拟登录成功，触发登录事件
			loginButton.SetValue(stateSession)
		}
	})

	log.Println("应用创建成功")
	log.Println("请在浏览器中访问 http://localhost:8507 查看应用")

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
