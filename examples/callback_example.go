package main

import (
"fmt"
"log"

"github.com/lengzhao/streamlit-go/app"
)

// CallbackExample 演示callback功能的示例
func CallbackExample(st *app.App) {
	// 添加标题
	st.Title("🔄 Callback功能演示")

	// 演示TextInputWithCallback
	st.Subheader("⌨️ TextInput With Callback")
	st.Text("输入文本时，会实时显示输入的值：")

	// 先创建一个WriteWidget来显示输入的姓名
	nameOutput := st.Write("")
	// 使用callback方式的TextInput
	nameInput := st.TextInputWithCallback("姓名:", func(value string) {
nameOutput.SetData("您输入的姓名是: " + value)
}, "")
	nameInput.SetPlaceholder("请输入姓名")

	// 演示NumberInputWithCallback
	st.Subheader("🔢 NumberInput With Callback")
	st.Text("输入数字时，会实时显示输入的值和平方：")

	// 先创建WriteWidgets来显示输入的数字和平方
	numberOutput := st.Write("")
	squareOutput := st.Write("")
	// 使用callback方式的NumberInput
	st.NumberInputWithCallback("数字:", func(value float64) {
numberOutput.SetData("您输入的数字是: " + fmt.Sprintf("%.0f", value))
squareOutput.SetData("该数字的平方是: " + fmt.Sprintf("%.0f", value*value))
}, 0)

	// 演示ButtonWithCallback
	st.Subheader("🔘 Button With Callback")
	st.Text("点击按钮时，会显示按钮被点击的消息：")

	// 先创建一个WriteWidget来显示按钮点击消息
	buttonOutput := st.Write("")
	// 使用callback方式的Button
	st.ButtonWithCallback("点击我!", func() {
		buttonOutput.SetData("按钮被点击了！")
	})
}

func main() {
	// 创建应用实例
	st := app.New(
app.WithTitle("Callback功能演示"),
app.WithPort(8502),
)

	// 运行Callback示例
	CallbackExample(st)

	log.Println("请在浏览器中访问 http://localhost:8502")

	// 启动应用
	if err := st.Run(); err != nil {
		log.Printf("服务器错误: %v", err)
	}
}
