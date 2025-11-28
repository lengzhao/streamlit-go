# 通信协议文档

## 1. 概述

Streamlit Go 使用 HTTP 进行客户端-服务器通信。HTTP GET 用于页面加载和静态资源传输，HTTP POST 用于组件事件处理和状态同步。

## 2. HTTP 通信

### 2.1 页面请求
- **路径**: `/`
- **方法**: GET
- **描述**: 获取应用主页面
- **响应**: HTML 页面，包含事件处理脚本

### 2.2 静态资源
- **路径**: `/static/*`
- **方法**: GET
- **描述**: 获取静态资源文件（CSS、JS等）
- **注意**: 当前版本未实现静态资源服务

### 2.3 健康检查
- **路径**: `/health`
- **方法**: GET
- **描述**: 服务器健康状态检查
- **响应**: `{"status": "ok"}`

### 2.4 组件事件
- **路径**: `/event`
- **方法**: POST
- **描述**: 处理组件事件（点击、输入等）
- **参数**:
  - `session_id`: 会话ID
  - `component_id`: 组件ID
  - `event_type`: 事件类型
  - `value`: 事件值

## 3. 消息格式

所有 HTTP POST 请求使用表单格式传递数据：

```
session_id=session_123456&component_id=widget_1&event_type=click&value=
```

HTTP 响应返回更新后的完整HTML内容：

```html
<div class="st-text" data-widget-id="widget_1">更新后的文本</div>
<button class="st-button" data-widget-id="widget_2">按钮</button>
```

## 4. 会话管理

### 4.1 会话ID生成
- 客户端首次访问时生成会话ID
- 格式: `session_{timestamp}_{random_string}`
- 会话ID存储在 Cookie 中并在URL参数中传递

### 4.2 会话关联
- HTTP请求通过 `sessionId` 参数与会话关联
- 服务端通过会话ID管理用户状态

### 4.3 会话超时
- 默认超时时间: 5分钟
- 清理间隔: 1分钟
- 超时后自动清理会话数据

## 5. 事件处理流程

1. 客户端触发事件（点击按钮、输入文本等）
2. JavaScript收集事件信息并发送HTTP POST请求到 `/event`
3. 服务端接收请求，查找对应组件并执行回调函数
4. 回调函数可能修改组件状态或会话数据
5. 服务端重新渲染所有组件并返回HTML
6. 客户端替换页面内容并重新绑定事件监听器

## 6. 安全考虑

### 6.1 CORS
- 当前版本未实现CORS限制
- 生产环境中应限制来源

### 6.2 会话安全
- 会话ID随机生成，难以猜测
- 会话数据隔离，用户间互不干扰

### 6.3 数据传输
- 建议在生产环境中使用 HTTPS
- 敏感数据应加密传输