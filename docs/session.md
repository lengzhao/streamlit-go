# 会话管理文档

## 1. 概述

Streamlit Go 提供了完善的会话管理机制，支持多用户隔离和状态管理。每个用户拥有独立的会话，确保数据安全和个性化体验。

## 2. 会话生命周期

### 2.1 会话创建
- 客户端首次访问时创建会话
- 会话ID通过Cookie持久化
- 会话数据存储在内存中

### 2.2 会话激活
- WebSocket连接建立时激活会话
- 会话最后访问时间更新

### 2.3 会话超时
- 默认超时时间：5分钟
- 定期清理过期会话（每1分钟）
- 超时后自动释放资源

### 2.4 会话销毁
- 手动删除会话
- 超时自动清理
- 服务器关闭时清理

## 3. 会话数据结构

### 3.1 Session 结构
```go
type Session struct {
    id             string                 // 会话唯一标识
    state          map[string]interface{} // 状态存储
    widgets        []widgets.Widget       // 会话私有组件
    widgetsMutex   sync.RWMutex           // 组件队列锁
    createdAt      time.Time              // 创建时间
    lastAccessedAt time.Time              // 最后访问时间
    mutex          sync.RWMutex           // 读写锁，保护并发访问
}
```

### 3.2 Manager 结构
```go
type Manager struct {
    sessions         map[string]*Session // 会话映射表
    mutex            sync.RWMutex        // 全局读写锁
    cleanupInterval  time.Duration       // 清理间隔
    sessionTimeout   time.Duration       // 会话超时时间
    cleanupCtx       context.Context     // 清理任务上下文
    cleanupCancel    context.CancelFunc  // 清理任务取消函数
    cleanupWaitGroup sync.WaitGroup      // 等待清理任务完成
}
```

## 4. 会话ID管理

### 4.1 ID生成策略
- 格式：`session_{timestamp}_{random_string}`
- 时间戳确保唯一性
- 随机字符串增加安全性

### 4.2 ID持久化
- 使用Cookie存储会话ID
- Cookie不设置过期时间（会话Cookie）
- 关闭浏览器后自动清除

### 4.3 ID传递
- 通过WebSocket连接参数传递
- URL编码确保传输安全

## 5. 状态管理

### 5.1 状态存储
- 使用键值对存储用户状态
- 支持任意类型的数据
- 线程安全的读写操作

### 5.2 状态操作
```go
// 设置状态
session.Set("key", value)

// 获取状态
value, exists := session.Get("key")

// 删除状态
session.Delete("key")

// 检查状态是否存在
exists := session.Has("key")

// 清空所有状态
session.Clear()
```

### 5.3 状态同步
- 状态变更自动更新最后访问时间
- 支持组件状态与会话状态联动

## 6. 组件会话绑定

### 6.1 全局组件
- 所有用户共享
- 存储在App的全局组件队列中
- 适用于公共信息展示

### 6.2 会话组件
- 每个用户独立拥有
- 存储在Session的组件队列中
- 适用于个性化内容

### 6.3 组件注册
```go
// 注册全局组件
app.AddWidget(widget)

// 注册会话组件
app.AddWidgetToSession(sessionID, widget)
```

## 7. 会话隔离

### 7.1 数据隔离
- 每个会话拥有独立的状态存储
- 组件数据互不干扰
- 确保用户隐私安全

### 7.2 组件隔离
- 会话组件只对特定用户可见
- 事件处理在正确的会话上下文中进行
- UI更新只影响目标用户

### 7.3 并发安全
- 使用读写锁保护会话数据
- 组件队列支持并发访问
- 避免竞态条件

## 8. 会话清理

### 8.1 自动清理
- 定期检查会话超时
- 异步清理过期会话
- 释放内存资源

### 8.2 手动清理
- 提供API删除指定会话
- 支持批量清理操作
- 灵活的会话管理

### 8.3 清理策略
- 基于最后访问时间
- 超时时间可配置
- 清理间隔可调整

## 9. 最佳实践

### 9.1 会话设计
- 合理划分全局组件和会话组件
- 避免在会话中存储大量数据
- 及时清理无用状态

### 9.2 性能优化
- 使用适当超时时间
- 避免频繁创建会话
- 合理设置清理间隔

### 9.3 安全考虑
- 保护会话ID不被泄露
- 验证会话ID的有效性
- 防止会话固定攻击