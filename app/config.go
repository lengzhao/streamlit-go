package app

// Config 应用配置
type Config struct {
	App struct {
		Title string // 应用标题
	} `json:"app"`
	Server struct {
		Host string `json:"host"` // 服务器主机
		Port int    `json:"port"` // 服务器端口
	} `json:"server"`
}

// Option 配置选项函数类型
type Option func(*Config)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	config := &Config{}
	config.App.Title = "Streamlit Go App"
	config.Server.Host = "localhost"
	config.Server.Port = 8501
	return config
}

// WithTitle 设置应用标题
func WithTitle(title string) Option {
	return func(c *Config) {
		c.App.Title = title
	}
}

// WithHost 设置服务器主机
func WithHost(host string) Option {
	return func(c *Config) {
		c.Server.Host = host
	}
}

// WithPort 设置服务器端口
func WithPort(port int) Option {
	return func(c *Config) {
		c.Server.Port = port
	}
}
