package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lance4117/gofuse/logger"
)

type Server struct {
	*gin.Engine
}

// New 初始化HTTP服务器引擎
func New(isDebug bool) *Server {
	// 颜色
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

	return &Server{engine}
}

// Run 启动HTTP服务器
func (s *Server) Run(addr ...string) {
	var path string
	// 默认本地8080
	if len(addr) == 0 {
		logger.Info("HTTP Service run at Default")
		path = "127.0.0.1:8080"
	} else {
		path = addr[0]
	}
	logger.Info("HTTP Service start at ", path)
	err := s.Engine.Run(path)
	if err != nil {
		logger.Fatal(err, "Server Run Fail")
		return
	}
}

// POST 注册POST请求处理函数
func (s *Server) POST(path string, handles ...ContextHandler) {
	s.Engine.POST(path, convertHandler(handles...)...)
}

// GET 注册GET请求处理函数
func (s *Server) GET(path string, handles ...ContextHandler) {
	s.Engine.GET(path, convertHandler(handles...)...)
}

// Use 注册中间件处理函数
func (s *Server) Use(handles ...ContextHandler) {
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
