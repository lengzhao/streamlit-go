package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lengzhao/streamlit-go/core"
	"github.com/lengzhao/streamlit-go/widgets"
)

func main() {
	// 创建服务实例
	service := core.NewService(
		core.WithTitle("会话Widgets示例"),
		core.WithPort(8504),
	)

	// 添加全局组件（所有用户共享）
	service.Title("🌐 全局标题 - 所有用户共享")
	service.Header("这是全局内容", true)
	service.Text("所有用户都会看到这段文字")

	// 为用户1创建按钮
	user1Button := widgets.NewButton("用户1按钮")

	user1Count := 0
	user1Button.OnChange(func(session widgets.ISession, event string, value string) {
		log.Println("Button clicked by user:", session.ID())
		user1Count++
		stat := widgets.NewText("用户计数器: " + fmt.Sprintf("%d", user1Count))
		session.AddWidget(stat)
	})
	service.AddWidget(user1Button)

	log.Println("服务创建成功")
	log.Println("请在浏览器中访问 http://localhost:8504 查看应用")

	// 设置信号处理，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 在单独的goroutine中启动服务
	go func() {
		if err := service.Start(); err != nil {
			log.Printf("服务器错误: %v", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("\n收到中断信号，关闭中...")

	// 优雅关闭
	if err := service.Stop(); err != nil {
		log.Printf("关闭时错误: %v", err)
	}

	log.Println("服务已成功停止")
}
