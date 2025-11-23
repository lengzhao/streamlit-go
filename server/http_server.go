package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/lengzhao/streamlit-go/ptemplate"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	host                         string
	port                         int
	hub                          *Hub
	server                       *http.Server
	eventHandler                 EventHandler
	getWidgetsForSessionCallback func(sessionID string) string // 为特定会话获取组件的回调
	appTitle                     string
}

// NewHTTPServer 创建新的HTTP服务器
func NewHTTPServer(host string, port int, hub *Hub) *HTTPServer {
	return &HTTPServer{
		host: host,
		port: port,
		hub:  hub,
	}
}

// SetEventHandler 设置事件处理器
func (s *HTTPServer) SetEventHandler(handler EventHandler) {
	s.eventHandler = handler
	// 同时设置Hub的事件处理器
	s.hub.SetEventHandler(handler)
}

// SetGetWidgetsForSessionCallback 设置为特定会话获取组件的回调函数
func (s *HTTPServer) SetGetWidgetsForSessionCallback(callback func(sessionID string) string) {
	s.getWidgetsForSessionCallback = callback
}

// SetAppTitle 设置应用标题
func (s *HTTPServer) SetAppTitle(title string) {
	s.appTitle = title
}

// Start 启动HTTP服务器
func (s *HTTPServer) Start() error {
	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.server = &http.Server{
		Addr: addr,
	}

	// 注册路由处理器
	s.registerRoutes()

	// 在单独的goroutine中启动Hub
	go s.hub.Run()

	// 启动HTTP服务器
	log.Printf("HTTP server starting on %s", addr)
	return s.server.ListenAndServe()
}

// Stop 停止HTTP服务器
func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// registerRoutes 注册路由处理器
func (s *HTTPServer) registerRoutes() {
	// 静态文件服务
	http.HandleFunc("/static/", s.serveStatic)

	// WebSocket连接
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWebSocket(s.hub, w, r)
	})

	// 主页
	http.HandleFunc("/", s.serveHome)

	// 健康检查
	http.HandleFunc("/health", s.serveHealth)
}

// serveStatic 处理静态文件请求
func (s *HTTPServer) serveStatic(w http.ResponseWriter, r *http.Request) {
	// 简单实现，实际项目中应该使用http.FileServer
	http.NotFound(w, r)
}

// serveHome 处理主页请求
func (s *HTTPServer) serveHome(w http.ResponseWriter, r *http.Request) {
	// 生成初始HTML页面
	html := s.generateInitialPage(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// serveHealth 处理健康检查请求
func (s *HTTPServer) serveHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// generateInitialPage 生成初始HTML页面
func (s *HTTPServer) generateInitialPage(r *http.Request) string {
	title := "Streamlit Go App"
	if s.appTitle != "" {
		title = s.appTitle
	}

	// 获取会话ID
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		// 如果没有提供会话ID，使用默认ID
		sessionID = "default-session-id"
	}

	// 生成组件HTML
	widgetsHTML := ""
	if s.getWidgetsForSessionCallback != nil {
		widgetsHTML = s.getWidgetsForSessionCallback(sessionID)
	}

	// 从模板包中获取页面模板
	tmpl, err := ptemplate.GetPageTemplate()
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return ""
	}

	data := map[string]interface{}{
		"Title":   title,
		"Content": template.HTML(widgetsHTML),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return ""
	}

	return buf.String()
}
