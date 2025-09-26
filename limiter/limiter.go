package limiter

import (
	"context"
	"net/http"
	"sync"

	"github.com/lance4117/gofuse/server"
	"golang.org/x/time/rate"
)

// Config 限流配置
type Config struct {
	Rate  float64 `yaml:"rate" json:"rate"`   // 每秒产生多少令牌
	Burst int     `yaml:"burst" json:"burst"` // 突发容量
}

// Manager 管理一组 limiter（支持按 key 区分，如 API 路径 / 用户ID）
type Manager struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	config   map[string]Config // 保存限流配置，便于动态调整
}

// NewLimiterManager 初始化
func NewLimiterManager(config map[string]Config) *Manager {
	return &Manager{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// GetLimiter 获取对应 key 的 limiter（不存在则创建）
func (m *Manager) GetLimiter(key string) *rate.Limiter {
	m.mu.RLock()
	lim, ok := m.limiters[key]
	m.mu.RUnlock()
	if ok {
		return lim
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	// double check
	if lim, ok := m.limiters[key]; ok {
		return lim
	}

	cfg := m.config[key]
	lim = rate.NewLimiter(rate.Limit(cfg.Rate), cfg.Burst)
	m.limiters[key] = lim
	return lim
}

// Allow 非阻塞尝试
func (m *Manager) Allow(key string) bool {
	return m.GetLimiter(key).Allow()
}

// Wait 阻塞等待（带 context）
func (m *Manager) Wait(ctx context.Context, key string) error {
	return m.GetLimiter(key).Wait(ctx)
}

// Middleware 返回中间件，按 API Path 进行限流
func Middleware(m *Manager) server.ContextHandler {
	return func(c *server.Context) {
		key := c.Path() // 根据路由路径作为 key

		if !m.Allow(key) {
			c.Fail(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
