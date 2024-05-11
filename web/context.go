package web

import (
	"github.com/caixr9527/go-cloud/web/render"
	"net/http"
	"sync"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	PathParams map[string]any // 路径参数
	Handlers   []Handler      // 中间件，控制器方法数组
	Index      int            // 指定路由的当前执行的方法索引
	StatusCode int            // 错误码
	Errors     error
	Size       int
	sameSite   http.SameSite
	Keys       map[string]any
	mu         sync.RWMutex
}

// 进入对应路由的下一个方法
func (c *Context) Next() {
	c.Index++
	for c.Index < len(c.Handlers) {
		c.Handlers[c.Index](c)
		c.Index++
	}
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) PathParam(key string) any {
	if v, ok := c.PathParams[key]; ok {
		return v
	}
	return nil
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.W.Header().Set(key, value)
}

func (c *Context) JSON(status int, data any) error {
	return c.render(status, &render.JSON{Data: data})
}

func (c *Context) XML(status int, data any) error {
	return c.render(status, &render.XML{Data: data})
}

func (c *Context) Redirect(status int, url string) error {
	return c.render(status, &render.Redirect{Code: status, Request: c.R, Location: url})
}

func (c *Context) render(statusCode int, r render.Render) error {
	err := r.Render(c.W, statusCode)
	c.StatusCode = statusCode
	return err
}
