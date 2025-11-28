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
	// 创建服务实例
	service := core.NewService(
		core.WithTitle("Widget 更新演示"),
		core.WithPort(8505),
	)

	// 添加全局组件（所有用户共享）
	service.Title("🌐 全局标题 - 所有用户共享")
	service.Header("这是全局内容", true)

	// 创建一个可更新的文本组件
	updateText := widgets.NewText("这是一个可以更新的文本")
	service.AddWidget(updateText)

	// 创建一个可删除的文本组件
	deleteText := widgets.NewText("这是一个可以删除的文本")
	service.AddWidget(deleteText)

	// 创建更新按钮
	updateButton := widgets.NewButton("更新文本")
	updateButton.OnChange(func(session widgets.ISession, event string, value string) {
		updateText.SetText("文本已更新！当前时间戳")
		session.SetWidget(updateText)
	})
	service.AddWidget(updateButton)

	// 创建删除按钮
	deleteButton := widgets.NewButton("删除文本")
	deleteButton.OnChange(func(session widgets.ISession, event string, value string) {
		session.DeleteWidget(deleteText.GetID())
	})
	service.AddWidget(deleteButton)

	log.Println("服务创建成功")
	log.Println("请在浏览器中访问 http://localhost:8505 查看应用")

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
