# 组件系统文档

## 1. 概述

Streamlit Go 的组件系统是框架的核心部分，提供了丰富的 UI 组件用于构建交互式应用。组件系统采用面向对象设计，具有良好的扩展性和可维护性。

## 2. 组件架构

### 2.1 Widget 接口
所有组件都必须实现 [Widget](file:///Volumes/ssd/myproject/streamlit-go/widgets/base.go#L27-L49) 接口：

```go
type Widget interface {
    // Render 渲染组件为HTML字符串
    Render() string

    // GetID 获取组件唯一标识
    GetID() string

    // GetType 获取组件类型
    GetType() string

    // SetKey 设置组件键值
    SetKey(key string)

    // GetKey 获取组件键值
    GetKey() string

    // OnChange 设置值变更回调函数
    OnChange(callback func(session SessionInterface, event string, value string))

    // IsVisible 检查组件是否可见
    IsVisible() bool
}
```

### 2.2 BaseWidget 基类
[BaseWidget](file:///Volumes/ssd/myproject/streamlit-go/widgets/base.go#L56-L63) 提供了所有组件的通用功能实现：

- 组件 ID 生成和管理
- 键值管理
- 回调函数管理
- 可见性控制
- UI 更新机制

### 2.3 ITriggerCallbacks 接口
[ITriggerCallbacks](file:///Volumes/ssd/myproject/streamlit-go/widgets/base.go#L51-L54) 接口用于触发组件回调：

```go
type ITriggerCallbacks interface {
    TriggerCallbacks(session SessionInterface, event string, value string)
}
```

## 3. 组件类型

### 3.1 文本组件
- Title: 标题组件
- Header: 头部组件
- Subheader: 子标题组件
- Text: 文本组件
- Write: 通用写入组件

### 3.2 输入组件
- TextInput: 文本输入组件
- NumberInput: 数字输入组件
- Button: 按钮组件

### 3.3 布局组件
- Container: 容器组件
- Columns: 列布局组件
- Sidebar: 侧边栏组件
- Expander: 可展开组件

### 3.4 数据展示组件
- Table: 表格组件
- DataFrame: 数据框组件
- Metric: 指标组件

## 4. 组件生命周期

### 4.1 创建
组件通过工厂函数创建，例如：
```go
button := widgets.NewButton("点击我")
```

### 4.2 注册
组件需要注册到应用中才能显示：
- 全局组件：`app.AddWidget(widget)`
- 会话组件：`app.AddWidgetToSession(sessionID, widget)`

### 4.3 渲染
组件通过 Render 方法生成 HTML：
```go
html := widget.Render()
```

### 4.4 事件处理
组件可以注册事件回调：
```go
button.OnChange(func(session widgets.SessionInterface, event string, value string) {
    // 处理事件
})
```

### 4.5 更新
组件状态变更后可以触发 UI 更新：
```go
widget.UpdateWidget(func() string {
    return widget.Render()
})
```

## 5. 会话组件 vs 全局组件

### 5.1 全局组件
- 所有用户共享
- 在 App 的 widgets 队列中管理
- 适用于全局信息展示

### 5.2 会话组件
- 每个用户独立拥有
- 在 Session 的 widgets 队列中管理
- 适用于用户个性化内容

## 6. 组件扩展

要创建自定义组件，需要：

1. 定义组件结构体，嵌入 BaseWidget：
```go
type CustomWidget struct {
    *BaseWidget
    // 自定义字段
}
```

2. 实现 Widget 接口的所有方法

3. 提供 Render 方法生成 HTML

4. 可选择实现 ITriggerCallbacks 接口