package web

import (
	"net/http"
	"sync"
)

type Context struct {
	W              http.ResponseWriter
	R              *http.Request
	Params         map[string]any
	GlobalHandlers []Handler // 全局中间件
	ReqHandlers    []Handler // 中间件，控制器方法数组
	Index          int       // 指定路由的当前执行的方法索引
	StatusCode     int       // 错误码
	Errors         error
	Size           int
	sameSite       http.SameSite
	Keys           map[string]any
	mu             sync.RWMutex
}

// 进入对应路由的下一个方法
func (c *Context) Next() {
	c.Index++
	for c.Index < len(c.ReqHandlers) {
		c.ReqHandlers[c.Index](c)
		c.Index++
	}
}
