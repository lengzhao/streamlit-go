package state

import (
"context"
"crypto/rand"
"encoding/hex"
"sync"
"time"
)

// Manager 状态管理器，管理所有会话
type Manager struct {
	sessions         map[string]*Session // 会话映射表
	mutex            sync.RWMutex        // 全局读写锁
	cleanupInterval  time.Duration       // 清理间隔
	sessionTimeout   time.Duration       // 会话超时时间
	cleanupCtx       context.Context     // 清理任务上下文
	cleanupCancel    context.CancelFunc  // 清理任务取消函数
	cleanupWaitGroup sync.WaitGroup      // 等待清理任务完成
}

// NewManager 创建新的状态管理器
func NewManager(cleanupInterval, sessionTimeout time.Duration) *Manager {
	return &Manager{
		sessions:        make(map[string]*Session),
		mutex:           sync.RWMutex{},
		cleanupInterval: cleanupInterval,
		sessionTimeout:  sessionTimeout,
	}
}

// GetSession 获取或创建会话
func (m *Manager) GetSession(sessionID string) *Session {
	m.mutex.RLock()
	session, exists := m.sessions[sessionID]
	m.mutex.RUnlock()

	if exists {
		return session
	}

	// 会话不存在，创建新会话
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 双重检查，防止并发创建
	session, exists = m.sessions[sessionID]
	if exists {
		return session
	}

	session = NewSession(sessionID)
	m.sessions[sessionID] = session
	return session
}

// DeleteSession 删除指定会话
func (m *Manager) DeleteSession(sessionID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.sessions, sessionID)
}

// CleanupExpiredSessions 清理过期会话
func (m *Manager) CleanupExpiredSessions() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for id, session := range m.sessions {
		if now.Sub(session.LastAccessedAt()) > m.sessionTimeout {
			delete(m.sessions, id)
		}
	}
}

// Start 启动定期清理任务
func (m *Manager) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cleanupCtx = ctx
	m.cleanupCancel = cancel

	m.cleanupWaitGroup.Add(1)
	go func() {
		defer m.cleanupWaitGroup.Done()
		ticker := time.NewTicker(m.cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.CleanupExpiredSessions()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop 停止清理任务
func (m *Manager) Stop() {
	if m.cleanupCancel != nil {
		m.cleanupCancel()
	}
	m.cleanupWaitGroup.Wait()
}

// SessionCount 返回当前会话数量
func (m *Manager) SessionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.sessions)
}

// GenerateSessionID 生成新的会话ID
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
