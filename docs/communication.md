# 通信协议文档

## 1. 概述

Streamlit Go 使用 HTTP 和 WebSocket 进行客户端-服务器通信。HTTP 用于页面加载和静态资源传输，WebSocket 用于实时交互和状态同步。

## 2. HTTP 通信

### 2.1 页面请求
- **路径**: `/`
- **方法**: GET
- **描述**: 获取应用主页面
- **响应**: HTML 页面，包含 WebSocket 连接脚本

### 2.2 静态资源
- **路径**: `/static/*`
- **方法**: GET
- **描述**: 获取静态资源文件（CSS、JS等）

### 2.3 健康检查
- **路径**: `/health`
- **方法**: GET
- **描述**: 服务器健康状态检查
- **响应**: `{"status": "ok"}`

## 3. WebSocket 通信

### 3.1 连接建立
- **路径**: `/ws`
- **参数**: `sessionId` - 会话ID
- **描述**: 建立 WebSocket 连接并关联会话

### 3.2 消息格式
所有 WebSocket 消息都使用 JSON 格式：

```json
{
  "type": "消息类型",
  "session_id": "会话ID",
  "data": "消息数据",
  "timestamp": "时间戳"
}
```

### 3.3 消息类型

#### 客户端发送的消息类型

1. **ping**
   - **描述**: 心跳消息
   - **数据**: 无
   - **示例**:
   ```json
   {
     "type": "ping",
     "session_id": "session_123456",
     "timestamp": 1234567890
   }
   ```

2. **component_event**
   - **描述**: 组件事件
   - **数据**: 
     - `componentId`: 组件ID
     - `eventType`: 事件类型（click, input等）
     - `value`: 事件值
   - **示例**:
   ```json
   {
     "type": "component_event",
     "session_id": "session_123456",
     "data": {
       "componentId": "widget_1",
       "eventType": "click",
       "value": ""
     },
     "timestamp": 1234567890
   }
   ```

#### 服务端发送的消息类型

1. **pong**
   - **描述**: 心跳响应
   - **数据**: 无
   - **示例**:
   ```json
   {
     "type": "pong",
     "session_id": "session_123456",
     "timestamp": 1234567890
   }
   ```

2. **ui_update**
   - **描述**: 完整 UI 更新
   - **数据**: 
     - `html`: 完整的 HTML 内容
   - **示例**:
   ```json
   {
     "type": "ui_update",
     "session_id": "session_123456",
     "data": {
       "html": "<div>...</div>"
     },
     "timestamp": 1234567890
   }
   ```

3. **partial_update**
   - **描述**: 局部 UI 更新
   - **数据**: 
     - `componentId`: 组件ID
     - `html`: 组件的 HTML 内容
   - **示例**:
   ```json
   {
     "type": "partial_update",
     "session_id": "session_123456",
     "data": {
       "componentId": "widget_1",
       "html": "<button>已点击</button>"
     },
     "timestamp": 1234567890
   }
   ```

4. **error**
   - **描述**: 错误消息
   - **数据**: 
     - `message`: 错误信息
     - `code`: 错误代码
   - **示例**:
   ```json
   {
     "type": "error",
     "session_id": "session_123456",
     "data": {
       "message": "组件未找到",
       "code": "COMPONENT_NOT_FOUND"
     },
     "timestamp": 1234567890
   }
   ```

## 4. 会话管理

### 4.1 会话ID生成
- 客户端首次访问时生成会话ID
- 格式: `session_{timestamp}_{random_string}`
- 会话ID存储在 Cookie 中

### 4.2 会话关联
- WebSocket 连接通过 `sessionId` 参数与会话关联
- 服务端通过会话ID管理用户状态

### 4.3 会话超时
- 默认超时时间: 5分钟
- 清理间隔: 1分钟
- 超时后自动清理会话数据

## 5. 心跳机制

### 5.1 客户端心跳
- 每30秒发送一次 ping 消息
- 检测连接状态

### 5.2 服务端心跳
- 每54秒发送一次 ping 消息
- 检测连接状态

### 5.3 连接重连
- 连接断开后自动重连
- 最多重连5次
- 重连间隔递增

## 6. 安全考虑

### 6.1 CORS
- 允许所有来源的 WebSocket 连接
- 生产环境中应限制来源

### 6.2 会话安全
- 会话ID随机生成，难以猜测
- 会话数据隔离，用户间互不干扰

### 6.3 数据传输
- 建议在生产环境中使用 HTTPS/WSS
- 敏感数据应加密传输