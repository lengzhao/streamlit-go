# 会话管理文档

## 1. 概述

Streamlit Go 使用会话(Session)来管理每个用户的独立状态。每个用户都有一个唯一的会话ID，用于区分不同用户的数据和组件状态。

## 2. 会话生命周期

### 2.1 会话创建
- 当用户首次访问应用时，系统会自动生成一个新的会话ID
- 会话ID通过Cookie和URL参数传递和保存
- 会话对象存储在State Manager中

### 2.2 会话使用
- 每次用户请求都会携带会话ID
- 服务端根据会话ID获取对应的会话对象
- 组件可以根据会话状态进行个性化渲染

### 2.3 会话销毁
- 会话默认超时时间为5分钟
- State Manager每1分钟清理一次过期会话
- 用户关闭浏览器标签页后，会话将在超时后自动清理

## 3. 会话ID管理

### 3.1 会话ID生成
会话ID采用以下格式生成：
```
session_{timestamp}_{random_string}
```

### 3.2 会话ID传递
会话ID通过以下方式在客户端和服务端之间传递：
1. URL参数：`?sessionId=session_123456`
2. Cookie：`streamlit_session_id=session_123456`

### 3.3 会话ID验证
- 服务端会验证会话ID的有效性
- 无效的会话ID会被重新生成

## 4. 会话数据存储

### 4.1 会话结构
每个会话包含以下数据：
- 会话ID
- 创建时间
- 最后访问时间
- 会话私有组件列表

### 4.2 组件存储
会话可以存储两种类型的组件：
- **全局组件**：所有用户共享的组件，存储在Service中
- **会话组件**：每个用户独立的组件，存储在Session中

## 5. 多用户支持

### 5.1 用户隔离
- 每个用户拥有独立的会话空间
- 用户之间的数据完全隔离
- 组件状态互不影响

### 5.2 并发访问
- 会话数据访问使用读写锁保护
- 支持多个用户同时访问
- 线程安全的组件操作

## 6. 会话API

### 6.1 Session接口
```go
type ISession interface {
    ID() string                           // 获取会话ID
    LastAccessedAt() time.Time            // 获取最后访问时间
    CreatedAt() time.Time                 // 获取创建时间
    AddWidget(widget Widget)              // 添加组件到会话
    SetWidget(widget Widget)              // 更新会话中的组件
    GetWidgets() []Widget                 // 获取会话组件列表
    ClearWidgets()                        // 清空会话组件
    DeleteWidget(componentID string)      // 从会话中删除组件
}
```

### 6.2 Manager接口
```go
type Manager struct {
    sessions        map[string]*Session   // 会话存储
    mutex           sync.RWMutex          // 读写锁
    cleanupInterval time.Duration         // 清理间隔
    sessionTimeout  time.Duration         // 会话超时时间
}

func (m *Manager) GetSession(sessionID string) *Session    // 获取或创建会话
func (m *Manager) DeleteSession(sessionID string)          // 删除会话
func (m *Manager) CleanupExpiredSessions()                // 清理过期会话
```

## 7. 使用示例

### 7.1 创建会话感知组件
```go
button := widgets.NewButton("点击计数")
count := 0
button.OnChange(func(session widgets.ISession, event string, value string) {
    count++
    // 为当前用户创建一个新的文本组件显示计数
    counter := widgets.NewText(fmt.Sprintf("计数: %d", count))
    session.AddWidget(counter)
})
```

### 7.2 会话数据操作
```go
// 在回调函数中访问会话数据
button.OnChange(func(session widgets.ISession, event string, value string) {
    // 添加组件到会话
    session.AddWidget(newWidget)
    
    // 删除会话中的组件
    session.DeleteWidget(componentID)
    
    // 获取会话信息
    sessionID := session.ID()
})
```

## 8. 最佳实践

### 8.1 会话组件设计
- 将用户特定的数据存储在会话组件中
- 全局共享的信息使用全局组件
- 避免在会话中存储大量数据

### 8.2 性能优化
- 合理设置会话超时时间
- 及时清理不必要的会话数据
- 使用组件复用减少内存占用

### 8.3 安全考虑
- 会话ID应足够随机，防止猜测攻击
- 敏感数据不应存储在客户端可访问的地方
- 定期清理过期会话释放资源