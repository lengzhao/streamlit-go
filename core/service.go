package core

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lengzhao/streamlit-go/ptemplate"
	"github.com/lengzhao/streamlit-go/state"
	"github.com/lengzhao/streamlit-go/widgets"
)

// Config 服务配置
type Config struct {
	Server struct {
		Host string
		Port int
	}
	App struct {
		Title string
	}
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Host string
			Port int
		}{
			Host: "localhost",
			Port: 8501,
		},
		App: struct {
			Title string
		}{
			Title: "Streamlit Go App",
		},
	}
}

// Service 核心服务
type Service struct {
	config        *Config
	stateManager  *state.Manager
	widgets       []widgets.Widget
	widgetsMutex  sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	server        *http.Server
	eventCallback func(session *state.Session, componentID string, eventType string, value string)
	callbackMutex sync.RWMutex
}

// Option 配置选项
type Option func(*Config)

// WithHost 设置主机地址
func WithHost(host string) Option {
	return func(c *Config) {
		c.Server.Host = host
	}
}

// WithPort 设置端口
func WithPort(port int) Option {
	return func(c *Config) {
		c.Server.Port = port
	}
}

// WithTitle 设置应用标题
func WithTitle(title string) Option {
	return func(c *Config) {
		c.App.Title = title
	}
}

// NewService 创建新的核心服务
func NewService(options ...Option) *Service {
	config := DefaultConfig()

	// 应用配置选项
	for _, opt := range options {
		opt(config)
	}

	// 创建状态管理器，会话超时5分钟，每1分钟清理一次
	stateManager := state.NewManager(1*time.Minute, 5*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())

	service := &Service{
		config:        config,
		stateManager:  stateManager,
		widgets:       make([]widgets.Widget, 0),
		ctx:           ctx,
		cancel:        cancel,
		eventCallback: nil,
	}

	// 设置全局更新函数
	widgets.SetGlobalUpdateFunc(func(componentID string, html string) {
		log.Printf("GlobalUpdateFunc: componentID=%s, html=%s", componentID, html)
	})

	return service
}

// Start 启动服务
func (s *Service) Start() error {
	log.Printf("Starting Streamlit Go service on %s:%d", s.config.Server.Host, s.config.Server.Port)

	// 启动状态管理器的清理任务
	s.stateManager.Start()

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.server = &http.Server{
		Addr: addr,
	}

	// 注册路由处理器
	s.registerRoutes()

	// 启动HTTP服务器
	log.Printf("HTTP server starting on %s", addr)
	return s.server.ListenAndServe()
}

// Stop 停止服务
func (s *Service) Stop() error {
	log.Println("Stopping Streamlit Go service")

	// 停止状态管理器
	s.stateManager.Stop()

	// 停止HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server stop error: %v", err)
		}
	}

	// 取消应用上下文
	if s.cancel != nil {
		s.cancel()
	}

	return nil
}

// AddWidget 添加组件到全局队列
func (s *Service) AddWidget(widget widgets.Widget) {
	s.widgetsMutex.Lock()
	defer s.widgetsMutex.Unlock()

	s.widgets = append(s.widgets, widget)
}

// GetWidgets 获取全局组件
func (s *Service) GetWidgets() []widgets.Widget {
	s.widgetsMutex.RLock()
	defer s.widgetsMutex.RUnlock()

	// 返回副本，避免并发问题
	widgetsCopy := make([]widgets.Widget, len(s.widgets))
	copy(widgetsCopy, s.widgets)
	return widgetsCopy
}

// Title 添加标题组件
func (s *Service) Title(text string) {
	title := widgets.NewText(text)
	title.SetKey("title")
	s.AddWidget(title)
}

// Header 添加头部组件
func (s *Service) Header(text string, withDivider bool) {
	header := widgets.NewText(text)
	header.SetKey("header")
	s.AddWidget(header)
}

// Text 添加文本组件
func (s *Service) Text(text string) {
	textWidget := widgets.NewText(text)
	s.AddWidget(textWidget)
}

// SetEventCallback 设置事件回调函数
func (s *Service) SetEventCallback(callback func(session *state.Session, componentID string, eventType string, value string)) {
	s.callbackMutex.Lock()
	defer s.callbackMutex.Unlock()
	s.eventCallback = callback
}

// handleEvent 处理事件
func (s *Service) handleEvent(session *state.Session, componentID string, eventType string, value string) {
	s.callbackMutex.RLock()
	defer s.callbackMutex.RUnlock()

	if s.eventCallback != nil {
		s.eventCallback(session, componentID, eventType, value)
	} else {
		// 默认处理：查找对应的组件并触发回调
		s.handleComponentEvent(session, componentID, eventType, value)
	}
}

// handleComponentEvent 处理组件事件
func (s *Service) handleComponentEvent(session *state.Session, componentID string, eventType string, value string) {
	log.Printf("Component event received: sessionID=%s, componentID=%s, eventType=%s, value=%v",
		session.ID(), componentID, eventType, value)

	// 查找对应的组件并更新其值（优先在会话组件中查找）
	var targetWidget widgets.Widget
	found := false

	// 先在会话组件中查找
	sessionWidgets := session.GetWidgets()
	for _, widget := range sessionWidgets {
		if widget.GetID() == componentID {
			targetWidget = widget
			found = true
			break
		}
	}

	// 如果没找到，再在全局组件中查找
	if !found {
		globalWidgets := s.GetWidgets()
		for _, widget := range globalWidgets {
			if widget.GetID() == componentID {
				targetWidget = widget
				found = true
				break
			}
		}
	}

	if found {
		log.Printf("Event widget: %s, Type: %s, Value: %v", targetWidget.GetID(), targetWidget.GetType(), value)
		bw, ok := targetWidget.(widgets.ITriggerCallbacks)
		if ok {
			// 在触发回调时传递会话对象，允许回调函数直接操作会话
			bw.TriggerCallbacks(session, eventType, value)
		} else {
			log.Printf("Widget does not implement ITriggerCallbacks")
		}
	} else {
		log.Printf("Widget with ID %s not found", componentID)
	}
}

// RenderWidgetsForSession 为指定会话渲染所有组件为HTML
func (s *Service) RenderWidgetsForSession(sessionID string) string {
	// 获取会话对象
	session := s.stateManager.GetSession(sessionID)

	// 获取全局组件
	globalWidgets := s.GetWidgets()

	// 获取指定会话组件
	sessionWidgets := session.GetWidgets()

	// 合并两个列表
	allWidgets := make([]widgets.Widget, 0, len(globalWidgets)+len(sessionWidgets))
	allWidgets = append(allWidgets, globalWidgets...)
	allWidgets = append(allWidgets, sessionWidgets...)

	// 直接渲染组件为HTML，不使用ui.Renderer
	html := ""
	for _, widget := range allWidgets {
		if widget.IsVisible() {
			html += widget.Render()
		}
	}

	return html
}

// GetAddress 获取服务器地址
func (s *Service) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
}

// registerRoutes 注册路由处理器
func (s *Service) registerRoutes() {
	// 静态文件服务
	http.HandleFunc("/static/", s.serveStatic)

	// WebSocket连接 - 简化版本中不实现
	http.HandleFunc("/ws", s.serveWebSocket)

	// 主页
	http.HandleFunc("/", s.serveHome)

	// 健康检查
	http.HandleFunc("/health", s.serveHealth)

	// 组件事件处理
	http.HandleFunc("/event", s.serveEvent)
}

// serveStatic 处理静态文件请求
func (s *Service) serveStatic(w http.ResponseWriter, r *http.Request) {
	// 简单实现，实际项目中应该使用http.FileServer
	http.NotFound(w, r)
}

// serveHome 处理主页请求
func (s *Service) serveHome(w http.ResponseWriter, r *http.Request) {
	// 生成初始HTML页面
	html := s.generateInitialPage(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// serveHealth 处理健康检查请求
func (s *Service) serveHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// serveEvent 处理组件事件
func (s *Service) serveEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单数据
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	sessionID := r.FormValue("session_id")
	componentID := r.FormValue("component_id")
	eventType := r.FormValue("event_type")
	value := r.FormValue("value")

	// 获取会话对象
	session := s.stateManager.GetSession(sessionID)
	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// 处理事件
	s.handleEvent(session, componentID, eventType, value)

	// 重新渲染页面并返回
	widgetsHTML := s.RenderWidgetsForSession(sessionID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(widgetsHTML))
}

// serveWebSocket 简化的WebSocket处理（不实现实际功能）
func (s *Service) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "WebSocket not implemented in simplified version", http.StatusNotImplemented)
}

// generateInitialPage 生成初始HTML页面
func (s *Service) generateInitialPage(r *http.Request) string {
	title := "Streamlit Go App"
	if s.config.App.Title != "" {
		title = s.config.App.Title
	}

	// 获取会话ID
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		// 生成新的会话ID
		var err error
		sessionID, err = state.GenerateSessionID()
		if err != nil {
			sessionID = "default-session-id"
		}
	}

	// 生成组件HTML
	widgetsHTML := s.RenderWidgetsForSession(sessionID)

	// 从模板包中获取页面模板
	tmpl, err := ptemplate.GetPageTemplate()
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return ""
	}

	data := map[string]interface{}{
		"Title":     title,
		"Content":   template.HTML(widgetsHTML),
		"SessionID": sessionID,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return ""
	}

	return buf.String()
}
