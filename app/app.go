package app

import (
"context"
"fmt"
"log"
"sync"
"time"

"github.com/lengzhao/streamlit-go/server"
"github.com/lengzhao/streamlit-go/state"
"github.com/lengzhao/streamlit-go/ui"
"github.com/lengzhao/streamlit-go/widgets"
)

// App 应用实例，管理整个应用的生命周期
type App struct {
	config              *Config                   // 应用配置
	stateManager        *state.Manager            // 状态管理器
	widgets             []widgets.Widget          // 全局组件队列（所有用户共享）
	sessionWidgets      map[string][]widgets.Widget // 每个会话的私有组件队列
	widgetsMutex        sync.RWMutex              // 全局组件队列锁
	sessionWidgetsMutex sync.RWMutex              // 会话组件映射锁
	currentSession      *state.Session            // 当前会话
	ctx                 context.Context           // 应用上下文
	cancel              context.CancelFunc        // 取消函数
	httpServer          *server.HTTPServer        // HTTP服务器
	hub                 *server.Hub               // WebSocket Hub
	renderer            *ui.Renderer              // UI渲染器
	currentSessionID    string                    // 当前会话ID
}

// New 创建新的应用实例
func New(options ...Option) *App {
	config := DefaultConfig()

	// 应用配置选项
	for _, opt := range options {
		opt(config)
	}

	// 创建状态管理器，会话超时30分钟，每5分钟清理一次
	stateManager := state.NewManager(5*time.Minute, 30*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())

	// 创建WebSocket Hub
	hub := server.NewHub()

	// 创建HTTP服务器
	httpServer := server.NewHTTPServer(config.Server.Host, config.Server.Port, hub)

	// 创建UI渲染器
	renderer := ui.NewRenderer()

	app := &App{
		config:              config,
		stateManager:        stateManager,
		widgets:             make([]widgets.Widget, 0),
		sessionWidgets:      make(map[string][]widgets.Widget),
		ctx:                 ctx,
		cancel:              cancel,
		httpServer:          httpServer,
		hub:                 hub,
		renderer:            renderer,
	}

	// 设置事件处理器
	app.httpServer.SetEventHandler(app)

	// 设置应用标题
	app.httpServer.SetAppTitle(config.App.Title)

	// 设置组件获取回调
	app.httpServer.SetGetWidgetsCallback(func() string {
		return app.RenderWidgets()
	})

	// 设置全局更新函数
	widgets.SetGlobalUpdateFunc(func(componentID string, html string) {
log.Printf("GlobalUpdateFunc: componentID=%s, html=%s", componentID, html)
if app.currentSessionID != "" {
if err := app.hub.SendPartialUpdate(app.currentSessionID, componentID, html); err != nil {
				log.Printf("Failed to send partial update: %v", err)
			}
		} else {
			log.Printf("GlobalUpdateFunc: currentSessionID is empty")
		}
	})

	return app
}

// Run 启动应用
func (a *App) Run() error {
	log.Printf("Starting Streamlit Go application on %s:%d", a.config.Server.Host, a.config.Server.Port)

	// 启动状态管理器的清理任务
	a.stateManager.Start()

	// 启动HTTP服务器
	return a.httpServer.Start()
}

// Stop 停止应用
func (a *App) Stop() error {
	log.Println("Stopping Streamlit Go application")

	// 停止HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.httpServer.Stop(ctx); err != nil {
		log.Printf("HTTP server stop error: %v", err)
	}

	// 停止状态管理器
	a.stateManager.Stop()

	// 取消应用上下文
	if a.cancel != nil {
		a.cancel()
	}

	return nil
}

// Rerun 重新运行应用（触发UI刷新）
func (a *App) Rerun() {
	// 清空全局组件队列
	a.widgetsMutex.Lock()
	a.widgets = make([]widgets.Widget, 0)
	a.widgetsMutex.Unlock()

	// 清空所有会话的组件队列
	a.sessionWidgetsMutex.Lock()
	a.sessionWidgets = make(map[string][]widgets.Widget)
	a.sessionWidgetsMutex.Unlock()

	// TODO: 触发应用脚本重新执行
}

// Session 获取当前会话
func (a *App) Session() *state.Session {
	if a.currentSession == nil {
		// 如果没有当前会话，创建一个默认会话
		sessionID, err := state.GenerateSessionID()
		if err != nil {
			log.Printf("Failed to generate session ID: %v", err)
			sessionID = "default"
		}
		a.currentSession = a.stateManager.GetSession(sessionID)
		a.currentSessionID = sessionID
	}
	return a.currentSession
}

// SetCurrentSession 设置当前会话
func (a *App) SetCurrentSession(sessionID string) {
	a.currentSession = a.stateManager.GetSession(sessionID)
	a.currentSessionID = sessionID
}

// AddWidget 添加组件到全局队列
func (a *App) AddWidget(widget widgets.Widget) {
	a.widgetsMutex.Lock()
	defer a.widgetsMutex.Unlock()

	a.widgets = append(a.widgets, widget)
}

// AddWidgetToSession 添加组件到指定会话的队列
func (a *App) AddWidgetToSession(sessionID string, widget widgets.Widget) {
	a.sessionWidgetsMutex.Lock()
	defer a.sessionWidgetsMutex.Unlock()

	if _, exists := a.sessionWidgets[sessionID]; !exists {
		a.sessionWidgets[sessionID] = make([]widgets.Widget, 0)
	}
	a.sessionWidgets[sessionID] = append(a.sessionWidgets[sessionID], widget)
}

// GetWidgets 获取全局组件
func (a *App) GetWidgets() []widgets.Widget {
	a.widgetsMutex.RLock()
	defer a.widgetsMutex.RUnlock()

	// 返回副本，避免并发问题
	widgetsCopy := make([]widgets.Widget, len(a.widgets))
	copy(widgetsCopy, a.widgets)
	return widgetsCopy
}

// GetSessionWidgets 获取指定会话的组件
func (a *App) GetSessionWidgets(sessionID string) []widgets.Widget {
	a.sessionWidgetsMutex.RLock()
	defer a.sessionWidgetsMutex.RUnlock()

	widgetsList, exists := a.sessionWidgets[sessionID]
	if !exists {
		return make([]widgets.Widget, 0)
	}

	// 返回副本，避免并发问题
	widgetsCopy := make([]widgets.Widget, len(widgetsList))
	copy(widgetsCopy, widgetsList)
	return widgetsCopy
}

// GetAllWidgets 获取所有组件（全局+当前会话）
func (a *App) GetAllWidgets() []widgets.Widget {
	// 获取全局组件
	globalWidgets := a.GetWidgets()

	// 获取当前会话组件
	var sessionWidgets []widgets.Widget
	if a.currentSessionID != "" {
		sessionWidgets = a.GetSessionWidgets(a.currentSessionID)
	}

	// 合并两个列表
	allWidgets := make([]widgets.Widget, 0, len(globalWidgets)+len(sessionWidgets))
	allWidgets = append(allWidgets, globalWidgets...)
	allWidgets = append(allWidgets, sessionWidgets...)

	return allWidgets
}

// ClearWidgets 清空全局组件队列
func (a *App) ClearWidgets() {
	a.widgetsMutex.Lock()
	defer a.widgetsMutex.Unlock()

	a.widgets = make([]widgets.Widget, 0)
}

// ClearSessionWidgets 清空指定会话的组件队列
func (a *App) ClearSessionWidgets(sessionID string) {
	a.sessionWidgetsMutex.Lock()
	defer a.sessionWidgetsMutex.Unlock()

	delete(a.sessionWidgets, sessionID)
}

// GetConfig 获取应用配置
func (a *App) GetConfig() *Config {
	return a.config
}

// GetStateManager 获取状态管理器
func (a *App) GetStateManager() *state.Manager {
	return a.stateManager
}

// GetAddress 获取服务器地址
func (a *App) GetAddress() string {
	return fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
}

// HandleComponentEvent 实现EventHandler接口，处理组件事件
func (a *App) HandleComponentEvent(sessionID string, event *server.ComponentEventData) {
	log.Printf("Component event received: sessionID=%s, componentID=%s, eventType=%s, value=%v",
sessionID, event.ComponentID, event.EventType, event.Value)

	// 设置当前会话ID，确保全局更新函数能正常工作
	a.SetCurrentSession(sessionID)

	// 查找对应的组件并更新其值（优先在会话组件中查找）
	var targetWidget widgets.Widget
	found := false

	// 先在会话组件中查找
	sessionWidgets := a.GetSessionWidgets(sessionID)
	for _, widget := range sessionWidgets {
		if widget.GetID() == event.ComponentID {
			targetWidget = widget
			found = true
			break
		}
	}

	// 如果没找到，再在全局组件中查找
	if !found {
		globalWidgets := a.GetWidgets()
		for _, widget := range globalWidgets {
			if widget.GetID() == event.ComponentID {
				targetWidget = widget
				found = true
				break
			}
		}
	}

	if found {
		log.Printf("Event widget: %s, Type: %s, Value: %v", targetWidget.GetID(), targetWidget.GetType(), event.Value)
		bw, ok := targetWidget.(widgets.ITriggerCallbacks)
		if ok {
			bw.TriggerCallbacks(event.EventType, event.Value)
		} else {
			log.Printf("Widget does not implement ITriggerCallbacks")
		}
	} else {
		log.Printf("Widget with ID %s not found", event.ComponentID)
	}
}

// RenderWidgets 渲染所有组件为HTML
func (a *App) RenderWidgets() string {
	widgetList := a.GetAllWidgets()
	// 添加调试日志
	log.Printf("RenderWidgets: rendering %d widgets", len(widgetList))
	for i, widget := range widgetList {
		log.Printf("RenderWidgets: widget %d, ID=%s, Type=%s", i, widget.GetID(), widget.GetType())
	}
	return a.renderer.RenderWidgets(widgetList)
}
