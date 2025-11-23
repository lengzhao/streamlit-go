# 部署指南

## 1. 环境要求

### 1.1 Go 版本
- Go 1.24.2 或更高版本
- 检查版本：`go version`

### 1.2 依赖库
- github.com/gorilla/websocket v1.5.3

### 1.3 系统要求
- Linux、macOS 或 Windows
- 至少 512MB 内存
- 至少 100MB 磁盘空间

## 2. 构建应用

### 2.1 获取代码
```bash
git clone https://github.com/lengzhao/streamlit-go.git
cd streamlit-go
```

### 2.2 安装依赖
```bash
go mod tidy
```

### 2.3 构建二进制文件
```bash
# 构建当前平台二进制
go build -o streamlit-go cmd/streamlit-go/main.go

# 构建 Linux 二进制
GOOS=linux GOARCH=amd64 go build -o streamlit-go-linux cmd/streamlit-go/main.go

# 构建 Windows 二进制
GOOS=windows GOARCH=amd64 go build -o streamlit-go.exe cmd/streamlit-go/main.go

# 构建 macOS 二进制
GOOS=darwin GOARCH=amd64 go build -o streamlit-go-mac cmd/streamlit-go/main.go
```

## 3. 部署方式

### 3.1 直接运行
```bash
# 运行示例应用
cd examples/basic
go run main.go

# 或运行构建的二进制文件
./streamlit-go
```

### 3.2 使用 systemd (Linux)
创建 systemd 服务文件 `/etc/systemd/system/streamlit-go.service`：

```ini
[Unit]
Description=Streamlit Go Application
After=network.target

[Service]
Type=simple
User=streamlit
WorkingDirectory=/path/to/streamlit-go
ExecStart=/path/to/streamlit-go/streamlit-go
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable streamlit-go
sudo systemctl start streamlit-go
```

### 3.3 使用 Docker
创建 Dockerfile：

```dockerfile
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o streamlit-go cmd/streamlit-go/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/streamlit-go .

EXPOSE 8501
CMD ["./streamlit-go"]
```

构建和运行：
```bash
# 构建镜像
docker build -t streamlit-go .

# 运行容器
docker run -p 8501:8501 streamlit-go
```

### 3.4 使用 Docker Compose
创建 docker-compose.yml：

```yaml
version: '3.8'

services:
  streamlit-go:
    build: .
    ports:
      - "8501:8501"
    volumes:
      - ./data:/root/data
    restart: unless-stopped
```

启动服务：
```bash
docker-compose up -d
```

## 4. 配置

### 4.1 环境变量
```bash
# 设置主机地址
export STREAMLIT_HOST=0.0.0.0

# 设置端口
export STREAMLIT_PORT=8501

# 设置应用标题
export STREAMLIT_TITLE="My Streamlit App"
```

### 4.2 命令行参数
在应用代码中支持命令行参数：

```go
import "flag"

func main() {
    var (
        host = flag.String("host", "localhost", "Host to bind")
        port = flag.Int("port", 8501, "Port to listen on")
        title = flag.String("title", "Streamlit Go App", "Application title")
    )
    
    flag.Parse()
    
    st := app.New(
        app.WithHost(*host),
        app.WithPort(*port),
        app.WithTitle(*title),
    )
    
    // ... 应用逻辑
}
```

### 4.3 配置文件
创建 config.json：

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8501
  },
  "app": {
    "title": "Streamlit Go App"
  }
}
```

在应用中读取配置：

```go
import (
    "encoding/json"
    "os"
)

type Config struct {
    Server struct {
        Host string `json:"host"`
        Port int    `json:"port"`
    } `json:"server"`
    App struct {
        Title string `json:"title"`
    } `json:"app"`
}

func loadConfig() (*Config, error) {
    file, err := os.Open("config.json")
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## 5. 反向代理

### 5.1 Nginx 配置
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8501;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /ws {
        proxy_pass http://localhost:8501;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5.2 Apache 配置
```apache
<VirtualHost *:80>
    ServerName your-domain.com
    
    ProxyPreserveHost On
    ProxyPass / http://localhost:8501/
    ProxyPassReverse / http://localhost:8501/
    
    # WebSocket support
    RewriteEngine On
    RewriteCond %{HTTP:Upgrade} websocket [NC]
    RewriteCond %{HTTP:Connection} upgrade [NC]
    RewriteRule ^/?(.*) "ws://localhost:8501/$1" [P,L]
</VirtualHost>
```

## 6. SSL/TLS 配置

### 6.1 使用 Let's Encrypt
```bash
# 安装 Certbot
sudo apt-get install certbot

# 获取证书
sudo certbot certonly --standalone -d your-domain.com

# 配置 Nginx 使用证书
ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
```

### 6.2 在应用中启用 HTTPS
```go
st := app.New(
    app.WithHost("0.0.0.0"),
    app.WithPort(8501),
)

// 在反向代理层处理 HTTPS
// 应用本身仍使用 HTTP
```

## 7. 监控和日志

### 7.1 日志配置
```go
import (
    "log"
    "os"
)

func setupLogging() {
    // 设置日志输出到文件
    file, err := os.OpenFile("streamlit-go.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }
    
    log.SetOutput(file)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```

### 7.2 健康检查
应用提供 `/health` 端点用于健康检查：
```bash
curl http://localhost:8501/health
# 返回: {"status": "ok"}
```

### 7.3 监控指标
```bash
# 检查进程状态
ps aux | grep streamlit-go

# 检查端口占用
netstat -tlnp | grep :8501

# 检查资源使用
top -p $(pgrep streamlit-go)
```

## 8. 故障排除

### 8.1 常见问题

#### 端口被占用
```bash
# 查找占用端口的进程
lsof -i :8501

# 杀死进程
kill -9 <PID>
```

#### 依赖问题
```bash
# 清理并重新下载依赖
go clean -modcache
go mod tidy
```

#### 权限问题
```bash
# 确保有足够的权限运行应用
chmod +x streamlit-go
sudo chown streamlit:streamlit /path/to/streamlit-go
```

### 8.2 调试技巧

#### 启用详细日志
```go
import "log"

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    // ... 应用逻辑
}
```

#### 检查会话状态
```bash
# 查看会话数量
curl http://localhost:8501/debug/sessions
```

## 9. 性能优化

### 9.1 资源限制
在 systemd 服务中设置资源限制：
```ini
[Service]
LimitNOFILE=65536
MemoryLimit=512M
```

### 9.2 连接数优化
调整系统文件描述符限制：
```bash
# 临时设置
ulimit -n 65536

# 永久设置（/etc/security/limits.conf）
* soft nofile 65536
* hard nofile 65536
```

### 9.3 会话超时配置
```go
// 在应用中配置会话超时
stateManager := state.NewManager(1*time.Minute, 5*time.Minute)
// 第一个参数是清理间隔，第二个参数是会话超时时间
```

## 10. 安全考虑

### 10.1 防火墙配置
```bash
# 只开放必要的端口
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable
```

### 10.2 用户权限
```bash
# 创建专用用户运行应用
sudo useradd -r -s /bin/false streamlit
sudo chown -R streamlit:streamlit /path/to/streamlit-go
```

### 10.3 定期更新
```bash
# 更新系统
sudo apt-get update && sudo apt-get upgrade

# 更新 Go 版本
# 从 https://golang.org/dl/ 下载最新版本
```