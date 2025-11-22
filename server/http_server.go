package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	host               string
	port               int
	hub                *Hub
	server             *http.Server
	eventHandler       EventHandler
	getWidgetsCallback func() string
	appTitle           string
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
}

// SetGetWidgetsCallback 设置获取组件回调函数
func (s *HTTPServer) SetGetWidgetsCallback(callback func() string) {
	s.getWidgetsCallback = callback
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
	html := s.generateInitialPage()
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
func (s *HTTPServer) generateInitialPage() string {
	title := "Streamlit Go App"
	if s.appTitle != "" {
		title = s.appTitle
	}

	// 生成组件HTML
	widgetsHTML := ""
	if s.getWidgetsCallback != nil {
		widgetsHTML = s.getWidgetsCallback()
	}

	// 生成会话ID（在实际应用中应该从请求中获取或生成）
	sessionID := "default-session-id" // 使用固定的会话ID以确保前后端一致

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            padding: 20px;
            margin: 0;
            background-color: #f5f5f5;
        }
        .st-container {
            max-width: 800px;
            margin: 0 auto;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
        }
        .st-title {
            color: #333;
            border-bottom: 1px solid #eee;
            padding-bottom: 10px;
            margin-bottom: 20px;
        }
        .st-header {
            color: #333;
            margin: 20px 0 10px 0;
        }
        .st-subheader {
            color: #666;
            margin: 15px 0 8px 0;
        }
        .st-text {
            color: #333;
            margin: 10px 0;
        }
        .st-button {
            background-color: #ff4b4b;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
            margin: 5px 0;
        }
        .st-button:hover {
            background-color: #ff3333;
        }
        .st-text-input-container, .st-number-input-container {
            margin: 10px 0;
        }
        .st-text-input-container label, .st-number-input-container label {
            display: block;
            margin-bottom: 5px;
            color: #333;
        }
        .st-text-input, .st-number-input {
            width: 100%%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
        }
        .st-container-with-border {
            border: 1px solid #ddd;
            padding: 15px;
            margin: 10px 0;
        }
        .st-columns {
            display: flex;
            gap: 20px;
            margin: 10px 0;
        }
        .st-column {
            flex: 1;
        }
        .st-sidebar {
            background-color: #f0f2f6;
            padding: 15px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .st-expander {
            border: 1px solid #ddd;
            border-radius: 4px;
            margin: 10px 0;
        }
        .st-expander-header {
            background-color: #f0f2f6;
            padding: 10px;
            cursor: pointer;
            font-weight: bold;
        }
        .st-expander-content {
            padding: 10px;
        }
        .st-table, .st-dataframe {
            width: 100%%;
            border-collapse: collapse;
            margin: 10px 0;
        }
        .st-table td, .st-dataframe td, .st-table th, .st-dataframe th {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        .st-table th, .st-dataframe th {
            background-color: #f0f2f6;
        }
        .st-metric {
            background-color: #f0f2f6;
            padding: 15px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .st-metric-label {
            font-size: 14px;
            color: #666;
            margin-bottom: 5px;
        }
        .st-metric-value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
        }
        .st-metric-delta {
            font-size: 14px;
            color: #666;
            margin-top: 5px;
        }
        .st-status {
            position: fixed;
            top: 10px;
            right: 10px;
            padding: 5px 10px;
            border-radius: 4px;
            font-size: 12px;
            background-color: #ff4b4b;
            color: white;
        }
    </style>
</head>
<body>
    <div class="st-container">
        <div id="connectionStatus" class="st-status">Connecting...</div>
        <div id="widgets-container">
            %s
        </div>
    </div>
    
    <script>
        let ws;
        let reconnectAttempts = 0;
        const maxReconnectAttempts = 5;
        
        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws?sessionId=' + encodeURIComponent("default-session-id");
            
            console.log('WebSocket connecting to:', wsUrl);
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                console.log('WebSocket connected');
                document.getElementById('connectionStatus').textContent = '✓ Connected';
                document.getElementById('connectionStatus').className = 'st-status';
                reconnectAttempts = 0;
            };
            
            ws.onmessage = function(event) {
                console.log('WebSocket message received:', event.data);
                handleMessage(JSON.parse(event.data));
            };
            
            ws.onclose = function() {
                console.log('WebSocket closed');
                document.getElementById('connectionStatus').textContent = '✗ Disconnected';
                document.getElementById('connectionStatus').className = 'st-status';
                
                // 尝试重连
                if (reconnectAttempts < maxReconnectAttempts) {
                    reconnectAttempts++;
                    setTimeout(connect, 1000 * reconnectAttempts);
                }
            };
            
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }
        
        function handleMessage(msg) {
            switch (msg.type) {
                case 'ui_update':
                    document.getElementById('widgets-container').innerHTML = msg.data.html;
                    break;
                case 'partial_update':
                    const element = document.querySelector('[data-widget-id="' + msg.data.componentId + '"]');
                    if (element) {
                        element.outerHTML = msg.data.html;
                    }
                    break;
                case 'error':
                    console.error('Server error:', msg.data.message);
                    break;
                default:
                    console.log('Unknown message type:', msg.type);
            }
        }
        
        function sendEvent(widgetId, eventType, value) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({
                    type: 'component_event',
                    sessionId: "default-session-id",
                    data: {
                        componentId: widgetId,
                        eventType: eventType,
                        value: value
                    },
                    timestamp: Date.now()
                }));
            }
        }
        
        // 页面加载完成后连接WebSocket
        window.addEventListener('load', function() {
            connect();
            
            // 为现有的元素添加事件监听器
            attachEventListeners();
        });
        
        function attachEventListeners() {
            // 按钮点击事件
            const buttons = document.querySelectorAll('[data-event-type="click"]');
            buttons.forEach(function(button) {
                button.addEventListener('click', function() {
                    sendEvent(this.dataset.widgetId, 'click', null);
                });
            });
            
            // 输入框变化事件
            const inputs = document.querySelectorAll('[data-event-type="input"]');
            inputs.forEach(function(input) {
                input.addEventListener('input', function() {
                    sendEvent(this.dataset.widgetId, 'input', this.value);
                });
            });
        }
        
        // 定期发送心跳
        setInterval(function() {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({
                    type: 'ping',
                    sessionId: "default-session-id",
                    timestamp: Date.now()
                }));
            }
        }, 30000);
    </script>
</body>
</html>
`, title, widgetsHTML, sessionID)
}
