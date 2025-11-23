# 示例代码目录

本目录包含了多个示例，演示了Streamlit Go的不同功能和用法。

## 目录结构

```
examples/
├── basic/              # 基础功能示例
│   └── main.go         # 基础组件使用示例
├── session-widgets/    # 会话Widgets示例
│   └── main.go         # 会话隔离功能示例
└── README.md           # 本文件
```

## 示例说明

### basic - 基础功能示例

展示了Streamlit Go的基础功能，包括：
- 各种文本组件（Title, Header, Subheader, Text, Write）
- 输入组件（TextInput, NumberInput, Button）
- 布局组件（Container, Columns, Sidebar, Expander）
- 数据展示组件（Table, DataFrame, Metric）

运行方式：
```bash
cd basic
go run main.go
```

访问 `http://localhost:8501` 查看示例应用。

### session-widgets - 会话Widgets示例

展示了Streamlit Go的会话隔离功能，包括：
- 全局组件（所有用户共享）
- 会话私有组件（每个用户独立）
- 会话状态隔离
- 多用户并发访问

运行方式：
```bash
cd session-widgets
go run main.go
```

访问以下URL进行测试：
- 用户1: `http://localhost:8503?sessionId=user-1-session`
- 用户2: `http://localhost:8503?sessionId=user-2-session`
- 默认用户: `http://localhost:8503?sessionId=default-session-id`

## 添加新的示例

要添加新的示例，请按照以下步骤操作：

1. 在examples目录下创建新的子目录
2. 在子目录中创建main.go文件
3. 实现相应的功能示例
4. 更新本README文件

例如：
```bash
mkdir -p examples/new_feature
# 在examples/new_feature/main.go中实现新功能示例
```