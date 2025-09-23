package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lance4117/gofuse/logger"
)

type Engine struct {
	*gin.Engine
}

// InitServer 初始化HTTP服务器引擎
func InitServer(isDebug bool) *Engine {
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

	return &Engine{engine}
}

// Run 启动HTTP服务器
func (e *Engine) Run(addr ...string) {
	var path string
	// 默认本地8080
	if len(addr) == 0 {
		logger.Info("HTTP Service run at Default")
		path = "127.0.0.1:8080"
	} else {
		path = addr[0]
	}
	logger.Info("HTTP Service start at ", path)
	err := e.Engine.Run(path)
	if err != nil {
		logger.Fatal(err, "Engine Run Fail")
		return
	}
}

// POST 注册POST请求处理函数
func (e *Engine) POST(path string, handles ...ContextHandler) {
	e.Engine.POST(path, convertHandler(handles...)...)
}

// GET 注册GET请求处理函数
func (e *Engine) GET(path string, handles ...ContextHandler) {
	e.Engine.GET(path, convertHandler(handles...)...)
}

// Use 注册中间件处理函数
func (e *Engine) Use(handles ...ContextHandler) {
	e.Engine.Use(convertHandler(handles...)...)
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
