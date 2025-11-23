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
	config              *Config                                                        // 应用配置
	stateManager        *state.Manager                                                 // 状态管理器
	widgets             []widgets.Widget                                               // 全局组件队列（所有用户共享）
	widgetsMutex        sync.RWMutex                                                   // 全局组件队列锁
	ctx                 context.Context                                                // 应用上下文
	cancel              context.CancelFunc                                             // 取消函数
	httpServer          *server.HTTPServer                                             // HTTP服务器
	hub                 *server.Hub                                                    // WebSocket Hub
	renderer            *ui.Renderer                                                   // UI渲染器
	loginCallback       func(session *state.Session)                                   // 登录回调函数
	appEventCallback    func(session *state.Session, event *server.ComponentEventData) // App级别事件回调函数
	appEventCallbackMux sync.RWMutex                                                   // App级别事件回调函数锁
}

// New 创建新的应用实例
func New(options ...Option) *App {
	config := DefaultConfig()

	// 应用配置选项
	for _, opt := range options {
		opt(config)
	}

	// 创建状态管理器，会话超时5分钟，每1分钟清理一次
	stateManager := state.NewManager(1*time.Minute, 5*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())

	// 创建WebSocket Hub
	hub := server.NewHub()

	// 设置状态管理器
	hub.SetStateManager(stateManager)
	// 同时设置StateManager的Hub引用，确保新创建的Session能获取到Hub
	stateManager.SetHub(hub)

	// 创建HTTP服务器
	httpServer := server.NewHTTPServer(config.Server.Host, config.Server.Port, hub)

	// 创建UI渲染器
	renderer := ui.NewRenderer()

	app := &App{
		config:           config,
		stateManager:     stateManager,
		widgets:          make([]widgets.Widget, 0),
		ctx:              ctx,
		cancel:           cancel,
		httpServer:       httpServer,
		hub:              hub,
		renderer:         renderer,
		appEventCallback: nil,
	}

	// 设置事件处理器
	app.httpServer.SetEventHandler(app)

	// 设置应用标题
	app.httpServer.SetAppTitle(config.App.Title)

	// 设置组件获取回调
	app.httpServer.SetGetWidgetsForSessionCallback(func(sessionID string) string {
		return app.RenderWidgetsForSession(sessionID)
	})

	// 设置全局更新函数
	widgets.SetGlobalUpdateFunc(func(componentID string, html string) {
		log.Printf("GlobalUpdateFunc: componentID=%s, html=%s", componentID, html)
		// 注意：由于移除了全局currentSession，这里需要通过其他方式确定目标会话
		// 在实际使用中，应该通过参数传递会话ID
		log.Printf("GlobalUpdateFunc: Warning - no current session context available")
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

	// TODO: 触发应用脚本重新执行
}

// AddWidget 添加组件到全局队列
func (a *App) AddWidget(widget widgets.Widget) {
	a.widgetsMutex.Lock()
	defer a.widgetsMutex.Unlock()

	a.widgets = append(a.widgets, widget)
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

// SetLoginCallback 设置登录回调函数
func (a *App) SetLoginCallback(callback func(session *state.Session)) {
	a.loginCallback = callback
}

// SetAppEventCallback 设置App级别事件回调函数
func (a *App) SetAppEventCallback(callback func(session *state.Session, event *server.ComponentEventData)) {
	a.appEventCallbackMux.Lock()
	defer a.appEventCallbackMux.Unlock()
	a.appEventCallback = callback
}

// handleAppLevelEvent 处理App级别的事件
func (a *App) handleAppLevelEvent(session *state.Session, event *server.ComponentEventData) {
	a.appEventCallbackMux.RLock()
	defer a.appEventCallbackMux.RUnlock()

	if a.appEventCallback != nil {
		a.appEventCallback(session, event)
	} else {
		log.Printf("No app-level event handler set for component ID %s", event.ComponentID)
	}
}

// HandleComponentEvent 实现EventHandler接口，处理组件事件
func (a *App) HandleComponentEvent(session *state.Session, event *server.ComponentEventData) {
	log.Printf("Component event received: sessionID=%s, componentID=%s, eventType=%s, value=%v",
		session.ID(), event.ComponentID, event.EventType, event.Value)

	// 查找对应的组件并更新其值（优先在会话组件中查找）
	var targetWidget widgets.Widget
	found := false

	// 先在会话组件中查找
	sessionWidgets := session.GetWidgets()
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
			// 在触发回调时传递会话对象，允许回调函数直接操作会话
			bw.TriggerCallbacks(session, event.EventType, event.Value)

			// 如果是登录事件，调用登录回调
			if event.EventType == "login" && a.loginCallback != nil {
				a.loginCallback(session)
			}
		} else {
			log.Printf("Widget does not implement ITriggerCallbacks")
		}
	} else {
		log.Printf("Widget with ID %s not found in session, trying app-level event handler", event.ComponentID)
		// 如果在Session中没找到widget，则调用app的事件回调接口
		a.handleAppLevelEvent(session, event)
	}
}

// RenderWidgetsForSession 为指定会话渲染所有组件为HTML
func (a *App) RenderWidgetsForSession(sessionID string) string {
	// 获取会话对象
	session := a.stateManager.GetSession(sessionID)

	// 获取全局组件
	globalWidgets := a.GetWidgets()

	// 获取指定会话组件
	sessionWidgets := session.GetWidgets()

	// 合并两个列表
	allWidgets := make([]widgets.Widget, 0, len(globalWidgets)+len(sessionWidgets))
	allWidgets = append(allWidgets, globalWidgets...)
	allWidgets = append(allWidgets, sessionWidgets...)

	// 添加调试日志
	log.Printf("RenderWidgetsForSession: rendering %d widgets for session %s", len(allWidgets), sessionID)
	for i, widget := range allWidgets {
		log.Printf("RenderWidgetsForSession: widget %d, ID=%s, Type=%s", i, widget.GetID(), widget.GetType())
	}
	return a.renderer.RenderWidgets(allWidgets)
}
