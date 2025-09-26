package limiter

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lance4117/gofuse/server"
)

func TestLimiter(t *testing.T) {
	// 假设配置从 YAML 加载
	cfg := map[string]Config{
		"/ping": {Rate: 0.1, Burst: 1},
	}

	lm := NewLimiterManager(cfg)

	s := server.New(true)
	s.Use(Middleware(lm))

	s.GET("/ping", func(ctx *server.Context) {
		ctx.OK(gin.H{"msg": "pong"})
	})

	s.Run(":8080")
}
