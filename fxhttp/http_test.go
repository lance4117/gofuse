package fxhttp

import (
	"testing"
)

func TestInitServer(t *testing.T) {
	engine := InitServer(true)
	engine.POST("/hello", Query1())

	engine.Run()
}

func Query1() ContextHandler {
	return func(ctx *Context) {

		ctx.Response(200, "hello world")
	}
}
