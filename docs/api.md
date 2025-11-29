# API 文档

## 1. App API

### 1.1 创建应用
```go
func New(options ...Option) *App
```
创建新的应用实例。

**参数**:
- `options`: 应用配置选项

**返回值**:
- `*App`: 应用实例

**示例**:
```go
app := app.New(
    app.WithTitle("我的应用"),
    app.WithPort(8080),
)
```

### 1.2 应用配置选项

#### WithTitle
```go
func WithTitle(title string) Option
```
设置应用标题。

#### WithHost
```go
func WithHost(host string) Option
```
设置应用主机地址。

#### WithPort
```go
func WithPort(port int) Option
```
设置应用端口。

### 1.3 应用控制

#### Run
```go
func (a *App) Run() error
```
启动应用。

#### Stop
```go
func (a *App) Stop() error
```
停止应用。

#### Rerun
```go
func (a *App) Rerun()
```
重新运行应用。

### 1.4 组件管理

#### AddWidget
```go
func (a *App) AddWidget(widget widgets.Widget)
```
添加全局组件。

#### AddWidgetToSession
```go
func (a *App) AddWidgetToSession(sessionID string, widget widgets.Widget)
```
添加会话组件。

#### GetWidgets
```go
func (a *App) GetWidgets() []widgets.Widget
```
获取全局组件列表。

### 1.5 回调设置

#### SetLoginCallback
```go
func (a *App) SetLoginCallback(callback func(session *state.Session))
```
设置登录回调函数。

#### SetAppEventCallback
```go
func (a *App) SetAppEventCallback(callback func(session *state.Session, event *server.ComponentEventData))
```
设置应用级事件回调函数。

### 1.6 获取器方法

#### GetConfig
```go
func (a *App) GetConfig() *Config
```
获取应用配置。

#### GetStateManager
```go
func (a *App) GetStateManager() *state.Manager
```
获取状态管理器。

#### GetHub
```go
func (a *App) GetHub() *server.Hub
```
获取WebSocket Hub。

#### GetAddress
```go
func (a *App) GetAddress() string
```
获取服务器地址。

## 2. 组件 API

### 2.1 文本组件

#### Title
```go
func (a *App) Title(text string) *widgets.TitleWidget
```
创建标题组件。

#### Header
```go
func (a *App) Header(text string, withDivider bool) *widgets.HeaderWidget
```
创建头部组件。

#### Subheader
```go
func (a *App) Subheader(text string) *widgets.SubheaderWidget
```
创建子标题组件。

#### Text
```go
func (a *App) Text(text string) *widgets.TextWidget
```
创建文本组件。

#### Write
```go
func (a *App) Write(data interface{}) *widgets.WriteWidget
```
创建写入组件。

### 2.2 输入组件

#### TextInput
```go
func (a *App) TextInput(label, value string) *widgets.TextInputWidget
```
创建文本输入组件。

#### NumberInput
```go
func (a *App) NumberInput(label string, value float64) *widgets.NumberInputWidget
```
创建数字输入组件。

#### Button
```go
func (a *App) Button(label string) *widgets.ButtonWidget
```
创建按钮组件。

### 2.3 布局组件

#### Container
```go
func (a *App) Container(withBorder bool) *widgets.ContainerWidget
```
创建容器组件。

#### Columns
```go
func (a *App) Columns(count int) *widgets.ColumnsWidget
```
创建列布局组件。

#### Sidebar
```go
func (a *App) Sidebar(key ...string) *widgets.SidebarWidget
```
创建侧边栏组件。

#### Expander
```go
func (a *App) Expander(label string, expanded bool) *widgets.ExpanderWidget
```
创建可展开组件。

### 2.4 数据展示组件

#### Table
```go
func (a *App) Table(data interface{}) *widgets.TableWidget
```
创建表格组件。

#### DataFrame
```go
func (a *App) DataFrame(data interface{}) *widgets.DataFrameWidget
```
创建数据框组件。

#### Metric
```go
func (a *App) Metric(label, value string) *widgets.MetricWidget
```
创建指标组件。

## 3. Session API

### 3.1 会话操作

#### GetSession
```go
func (m *Manager) GetSession(sessionID string) *Session
```
获取或创建会话。

#### DeleteSession
```go
func (m *Manager) DeleteSession(sessionID string)
```
删除指定会话。

#### CleanupExpiredSessions
```go
func (m *Manager) CleanupExpiredSessions()
```
清理过期会话。

### 3.2 会话状态管理

#### Set
```go
func (s *Session) Set(key string, value interface{})
```
设置会话状态。

#### Get
```go
func (s *Session) Get(key string) (interface{}, bool)
```
获取会话状态。

#### Delete
```go
func (s *Session) Delete(key string)
```
删除会话状态。

#### Has
```go
func (s *Session) Has(key string) bool
```
检查状态是否存在。

#### Clear
```go
func (s *Session) Clear()
```
清空所有状态。

### 3.3 组件管理

#### AddWidget
```go
func (s *Session) AddWidget(widget widgets.Widget)
```
添加组件到会话。

#### GetWidgets
```go
func (s *Session) GetWidgets() []widgets.Widget
```
获取会话组件列表。

#### ClearWidgets
```go
func (s *Session) ClearWidgets()
```
清空会话组件。

## 4. Widget API

### 4.1 通用方法

#### Render
```go
func (w *BaseWidget) Render() string
```
渲染组件为HTML字符串。

#### GetID
```go
func (w *BaseWidget) GetID() string
```
获取组件唯一标识。

#### GetType
```go
func (w *BaseWidget) GetType() string
```
获取组件类型。

#### SetKey
```go
func (w *BaseWidget) SetKey(key string)
```
设置组件键值。

#### GetKey
```go
func (w *BaseWidget) GetKey() string
```
获取组件键值。

#### OnChange
```go
func (w *BaseWidget) OnChange(callback func(session ISession, event string, value string))
```
设置值变更回调函数。

#### IsVisible
```go
func (w *BaseWidget) IsVisible() bool
```
检查组件是否可见。

### 4.2 特定组件方法

#### TextInput.SetValue
```go
func (w *TextInputWidget) SetValue(value string)
```
设置文本输入值。

#### NumberInput.SetValue
```go
func (w *NumberInputWidget) SetValue(value float64)
```
设置数字输入值。

#### Button.SetValue
```go
func (w *ButtonWidget) SetValue(session ISession)
```
触发按钮点击事件。

#### Write.SetData
```go
func (w *WriteWidget) SetData(data interface{})
```
设置写入组件数据。

#### Metric.SetDelta
```go
func (w *MetricWidget) SetDelta(delta string)
```
设置指标变化值。

## 5. Server API

### 5.1 Hub 方法

#### NewHub
```go
func NewHub() *Hub
```
创建新的Hub。

#### Run
```go
func (h *Hub) Run()
```
运行Hub。

#### Register
```go
func (h *Hub) Register(client *Client)
```
注册客户端。

#### Unregister
```go
func (h *Hub) Unregister(client *Client)
```
注销客户端。

#### Broadcast
```go
func (h *Hub) Broadcast(message []byte)
```
广播消息。

#### SendToSession
```go
func (h *Hub) SendToSession(sessionID string, message []byte)
```
发送消息到指定会话。

#### SendUIUpdate
```go
func (h *Hub) SendUIUpdate(sessionID string, html string) error
```
发送UI更新到指定会话。

#### SendPartialUpdate
```go
func (h *Hub) SendPartialUpdate(sessionID string, componentID string, html string) error
```
发送局部更新到指定会话。

#### SendError
```go
func (h *Hub) SendError(sessionID string, errorMsg string, errorCode string) error
```
发送错误消息到指定会话。

### 5.2 HTTP 服务器方法

#### NewHTTPServer
```go
func NewHTTPServer(host string, port int, hub *Hub) *HTTPServer
```
创建新的HTTP服务器。

#### Start
```go
func (s *HTTPServer) Start() error
```
启动HTTP服务器。

#### Stop
```go
func (s *HTTPServer) Stop(ctx context.Context) error
```
停止HTTP服务器。

#### SetEventHandler
```go
func (s *HTTPServer) SetEventHandler(handler server.EventHandler)
```
设置事件处理器。

#### SetAppTitle
```go
func (s *HTTPServer) SetAppTitle(title string)
```
设置应用标题。

#### SetGetWidgetsForSessionCallback
```go
func (s *HTTPServer) SetGetWidgetsForSessionCallback(callback func(sessionID string) string)
```
设置获取会话组件回调函数。