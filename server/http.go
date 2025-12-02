package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lance4117/gofuse/logger"
)

type HttpServer struct {
	*gin.Engine
}

// NewHTTP 初始化HTTP服务引擎
func NewHTTP(isDebug bool) *HttpServer {
	// 彩色
	gin.ForceConsoleColor()
	// 是否为debug模式
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 获取engine
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	return &HttpServer{engine}
}

// Run 启动HTTP服务，返回错误交由上层决策，避免库代码强退进程。
func (s *HttpServer) Run(addr ...string) error {
	var path string
	// 默认监听8080
	if len(addr) == 0 {
		logger.Info("HTTP Service run at Default")
		path = "127.0.0.1:8080"
	} else {
		path = addr[0]
	}
	logger.Info("HTTP Service start at ", path)
	if err := s.Engine.Run(path); err != nil {
		return err
	}
	return nil
}

// POST 注册POST路由
func (s *HttpServer) POST(path string, handles ...ContextHandler) {
	s.Engine.POST(path, convertHandler(handles...)...)
}

// GET 注册GET路由
func (s *HttpServer) GET(path string, handles ...ContextHandler) {
	s.Engine.GET(path, convertHandler(handles...)...)
}

// Use 注册中间件
func (s *HttpServer) Use(handles ...ContextHandler) {
	s.Engine.Use(convertHandler(handles...)...)
}

// convertHandler 将ContextHandler转换为gin.HandlerFunc
func convertHandler(handlers ...ContextHandler) []gin.HandlerFunc {
	var ginHandlers []gin.HandlerFunc
	for _, handler := range handlers {
		h := handler
		ginHandlers = append(ginHandlers, func(ctx *gin.Context) {
			h(&Context{ctx})
		})
	}
	return ginHandlers
}
