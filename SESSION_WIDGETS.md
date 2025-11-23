# 会话Widgets功能说明

## 功能概述

本项目现已支持为不同用户创建独立的Widgets，这些Widgets与用户的会话(session)挂钩。通过这种方式，可以实现：

1. 全局组件：所有用户共享的组件
2. 用户私有组件：仅特定用户可见的组件
3. 会话隔离：不同用户之间的组件状态互不干扰

## 实现细节

### 核心改动

1. 在[Session](file:///Volumes/ssd/myproject/streamlit-go/state/session.go#L8-L15)结构中直接添加了会话组件管理：
   - [widgets](file:///Volumes/ssd/myproject/streamlit-go/state/session.go#L11-L11): 存储会话的私有组件队列
   - [widgetsMutex](file:///Volumes/ssd/myproject/streamlit-go/state/session.go#L12-L12): 保护组件队列的读写锁
   - [AddWidget](file:///Volumes/ssd/myproject/streamlit-go/state/session.go#L99-L105): 添加组件到会话
   - [GetWidgets](file:///Volumes/ssd/myproject/streamlit-go/state/session.go#L108-L116): 获取会话组件
   - [ClearWidgets](file:///Volumes/ssd/myproject/streamlit-go/state/session.go#L119-L125): 清空会话组件

2. 从[App](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L17-L31)结构中移除了全局[currentSession](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L24-L24)状态，确保会话隔离：
   - 移除了[currentSession](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L24-L24)字段
   - 移除了[currentSessionID](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L30-L30)字段
   - 移除了[Session()](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L142-L150)和[SetCurrentSession()](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L153-L156)方法

3. 修改了会话组件管理方法：
   - [GetSessionWidgets](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L189-L203): 从Session对象中获取组件
   - [ClearSessionWidgets](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L232-L238): 清空Session对象中的组件
   - [GetAllWidgets](file:///Volumes/ssd/myproject/streamlit-go/app/app.go#L205-L222): 获取全局组件和指定会话的组件

4. 为所有Widget API方法保留了会话版本以保持向后兼容：
   - `Title` → `TitleToSession` / `TitleToCurrentSession`
   - `Text` → `TextToSession` / `TextToCurrentSession`
   - `Button` → `ButtonToSession` / `ButtonToCurrentSession`
   - 以及其他所有组件类型

### 架构优势

1. **真正的会话隔离**：每个用户的组件状态完全独立存储在自己的Session对象中
2. **线程安全**：通过为每个Session对象添加读写锁，确保并发访问安全
3. **清晰的数据结构**：会话数据和组件直接存储在Session对象中，而非通过外部映射管理
4. **更好的性能**：直接访问Session对象中的组件，避免了映射查找的开销

### 使用方法

#### 1. 为特定会话添加组件

```go
// 为用户1添加组件
user1SessionID := "user-1-session"
st.TitleToSession(user1SessionID, "用户1的专属标题")

// 为用户2添加组件
user2SessionID := "user-2-session"
st.TextToSession(user2SessionID, "用户2的专属文本")
```

#### 2. 为当前会话添加组件

```go
// 在处理WebSocket事件时，可以通过参数传递会话ID
st.TitleToSession(currentSessionID, "当前用户的标题")
```

#### 3. 处理会话组件事件

```go
user1Button := st.ButtonToSession(user1SessionID, "用户1按钮")
user1Button.OnChange(func(event string, value string) {
    // 当按钮被点击时，更新用户1的输出
    output := st.WriteToSession(user1SessionID, "")
    output.SetData("用户1按钮被点击了！")
})
```

## 技术实现要点

1. **会话ID传递**：通过WebSocket连接的查询参数传递会话ID
2. **组件查找**：在处理组件事件时，优先在会话组件中查找，找不到再查找全局组件
3. **状态隔离**：每个会话的组件状态独立存储，互不干扰
4. **渲染合并**：渲染时合并全局组件和当前会话组件

## 会话状态管理

为了确保不同用户之间的状态完全隔离，需要注意以下几点：

1. **组件创建时机**：用户特定的组件应在会话初始化时创建，而不是在应用启动时全局创建
2. **状态独立性**：每个用户的输入、按钮点击等交互状态都应独立存储
3. **页面刷新处理**：刷新页面时应保持当前用户的状态，而不是显示其他用户的状态

## 使用示例

在实际应用中，会话ID通过WebSocket连接自动传递：

```javascript
// 前端连接示例
const ws = new WebSocket(`ws://localhost:8501/ws?sessionId=${userSessionId}`);
```

后端会自动将该连接与指定会话关联，后续的组件操作都会针对该会话进行。

### 用户会话隔离示例

```go
// 为当前会话创建用户特定的组件
st.HeaderToSession(userSessionID, "👤 您的个人空间", true)

// 创建用户特定的输入组件
nameOutput := st.WriteToSession(userSessionID, "")
nameInput := st.TextInputToSession(userSessionID, "姓名", "")
nameInput.OnChange(func(name string) {
    if name != "" {
        nameOutput.SetData("您好，" + name + "！")
    }
})

// 创建用户特定的按钮
buttonOutput := st.WriteToSession(userSessionID, "")
button := st.ButtonToSession(userSessionID, "点击我")
button.OnChange(func() {
    buttonOutput.SetData(fmt.Sprintf("按钮被点击了！时间：%s", time.Now().Format("2006-01-02 15:04:05")))
})
```

通过这种方式，每个用户都有自己独立的组件实例和状态，不同用户之间的操作不会相互影响。