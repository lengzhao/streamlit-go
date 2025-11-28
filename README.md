# Streamlit Go

Streamlit Go 是一个用 Go 语言实现的 Streamlit 框架，用于快速构建交互式数据应用。

## 功能特性

- **基础框架**：应用结构、配置管理
- **状态管理**：Session会话与StateManager，支持自动过期
- **通信机制**：基于HTTP POST请求的组件事件处理和页面更新
- **组件系统**：支持文本（Title/Text）、输入（TextInput/Button）等组件，具备注册机制与ID生成
- **UI渲染**：集成HTML模板、CSS样式与JavaScript脚本
- **HTTP服务**：基于net/http提供路由、静态资源与健康检查
- **会话隔离Widgets**：支持为不同用户创建独立的Widgets，实现多用户状态隔离

## 安装

```bash
go mod tidy
```

## 运行示例

请查看 [examples/README.md](./examples/README.md) 了解不同示例的运行方式。

### 会话Widgets示例

新增的会话Widgets示例演示了如何为不同用户创建独立的界面组件：

```bash
cd examples/session-widgets
go run main.go
```

然后在浏览器中访问以下URL进行测试：
- 用户1: http://localhost:8504?sessionId=user-1-session
- 用户2: http://localhost:8504?sessionId=user-2-session
- 默认用户: http://localhost:8504?sessionId=default-session-id

### Widget更新演示示例

Widget更新演示示例展示了如何动态更新和删除组件：

```bash
cd examples/widget-update-demo
go run main.go
```

然后在浏览器中访问 http://localhost:8505 进行测试。

## 会话Widgets功能

本项目支持为不同用户创建独立的Widgets，这些Widgets与用户的会话(session)挂钩。

### 会话隔离特性

- 每个用户都有独立的会话状态
- 不同浏览器窗口代表不同用户
- 每个用户的输入和按钮点击都是独立的
- 刷新页面会保持当前用户的状态

## 目录结构

```
.
├── core/        # 核心服务实现
├── examples/    # 示例代码
├── ptemplate/   # 页面模板
├── state/       # 状态与会话管理
├── widgets/     # 所有UI组件实现
├── go.mod       # Go模块定义
└── go.sum       # Go模块校验和
```

## 组件类型

- 文本组件：Title, Header, Subheader, Text, Write
- 输入组件：TextInput, NumberInput, Button
- 布局组件：Container, Columns, Sidebar, Expander
- 数据展示：Table, DataFrame, Metric

## 许可证

MIT