package limiter

import (
	"context"
	"net/http"
	"sync"

	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/server"
	"golang.org/x/time/rate"
)

// Config 速率配置
type Config struct {
	Rate  float64 `yaml:"rate" json:"rate"`   // 每秒允许的请求数
	Burst int     `yaml:"burst" json:"burst"` // 峰值令牌桶容量
}

// Manager 管理一组 limiter，支持按 key 区分，如 API 路径 / 用户ID 等
type Manager struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	config   map[string]Config // 配置表，可动态更新
}

// NewLimiterManager 初始化
func NewLimiterManager(config map[string]Config) *Manager {
	return &Manager{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// GetLimiter 获取对应 key 的 limiter，不存在则创建
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

	cfg := m.effectiveConfig(key)
	lim = rate.NewLimiter(rate.Limit(cfg.Rate), cfg.Burst)
	m.limiters[key] = lim
	return lim
}

func (m *Manager) effectiveConfig(key string) Config {
	cfg, ok := m.config[key]
	if !ok {
		cfg = Config{Rate: 5, Burst: 5} // 全局兜底，避免零值导致拒绝所有请求
	}
	if cfg.Rate <= 0 {
		cfg.Rate = 5
	}
	if cfg.Burst <= 0 {
		cfg.Burst = 5
	}
	return cfg
}

// Allow 尝试放行一次
func (m *Manager) Allow(key string) bool {
	return m.GetLimiter(key).Allow()
}

// Wait 阻塞等待令牌，支持 context
func (m *Manager) Wait(ctx context.Context, key string) error {
	return m.GetLimiter(key).Wait(ctx)
}

// Middleware 基于 API Path 的限流中间件
func Middleware(m *Manager) server.ContextHandler {
	return func(c *server.Context) {
		key := c.Path() // 使用路由路径作为 key

		if !m.Allow(key) {
			logger.Errorf("%s has too many requests", key)
			c.Fail(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
