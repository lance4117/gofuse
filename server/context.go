package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Context struct {
	GinCtx *gin.Context
}

type Account struct {
	Uid     int64
	AppId   int64
	Lang    string
	Version string
	Os      int
	Ltm     int64
}

type ContextHandler func(ctx *Context)

// Account 返回用户账号相关信息
func (c *Context) Account() *Account {
	v, _ := c.GinCtx.Get("account")
	a, ok := v.(*Account)
	if ok {
		return a
	}
	return &Account{}
}

// SetAccount 设置用户账号相关信息
func (c *Context) SetAccount(account *Account) {
	if account == nil {
		return
	}
	c.GinCtx.Set("account", *account)
}

// Header 返回http请求头
func (c *Context) Header() http.Header {
	return c.GinCtx.Request.Header
}

// Next 执行下一个handler
func (c *Context) Next() {
	c.GinCtx.Next()
}

// Abort 终端执行
func (c *Context) Abort() {
	c.GinCtx.Abort()
}

// OK 响应成功，data为响应的数据
func (c *Context) OK(data any) {
	c.Response(http.StatusOK, data)
}

// Fail 响应失败，data为响应的数据
func (c *Context) Fail(code int) {
	c.Response(code, nil)
}

// Response 响应请求，code为状态码,data 响应的数据
func (c *Context) Response(code int, data any) {
	c.GinCtx.Header("Code", strconv.Itoa(code))
	switch c.ContentType() {
	case binding.MIMEJSON:
		c.GinCtx.JSON(code, data)
	case binding.MIMEPROTOBUF:
		c.GinCtx.ProtoBuf(code, data)
	default:
		c.GinCtx.JSON(code, data)
	}
}

// IP 返回客户端真实IP
func (c *Context) IP() string {
	return c.GinCtx.ClientIP()
}

// Path 返回客户端完整地址
func (c *Context) Path() string {
	return c.GinCtx.FullPath()
}

// ContentType 返回 Content-Type header
func (c *Context) ContentType() string {
	return c.GinCtx.ContentType()
}

// Bind 绑定请求体到obj
func (c *Context) Bind(obj any) error {
	return c.GinCtx.ShouldBind(obj)
}

// BindHeader 绑定header
func (c *Context) BindHeader(header any) error {
	return c.GinCtx.BindHeader(header)
}

// BindUri 绑定uri
func (c *Context) BindUri(uri any) error {
	return c.GinCtx.BindUri(uri)
}
