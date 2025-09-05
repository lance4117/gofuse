package fshttp

import (
	"gitee.com/lance4117/GoFuse/fslogger"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
}

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

func (e *Engine) Run(addr ...string) {
	var path string
	// 默认本地8080
	if len(addr) == 0 {
		fslogger.Info("HTTP Service run at Default")
		path = "127.0.0.1:8080"
	} else {
		path = addr[0]
	}
	fslogger.Info("HTTP Service start at ", path)
	err := e.Engine.Run(path)
	if err != nil {
		fslogger.Fatal("Engine Run Fail")
		return
	}
}

func (e *Engine) POST(path string, handles ...ContextHandler) {
	e.Engine.POST(path, convertHandler(handles...)...)
}

func (e *Engine) GET(path string, handles ...ContextHandler) {
	e.Engine.GET(path, convertHandler(handles...)...)
}

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
