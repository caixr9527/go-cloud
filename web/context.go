package web

import (
	"errors"
	"github.com/caixr9527/go-cloud/common/utils/stringUtils"
	"github.com/caixr9527/go-cloud/factory"
	logger "github.com/caixr9527/go-cloud/log"
	"github.com/caixr9527/go-cloud/web/binding"
	"github.com/caixr9527/go-cloud/web/render"
	"html/template"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const defaultMaxMemory = 32 << 20

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	Params     map[string]any // 路径参数
	Handlers   []Handler      // 中间件，控制器方法数组
	Index      int            // 指定路由的当前执行的方法索引
	StatusCode int            // 错误码
	Errors     error
	Size       int
	sameSite   http.SameSite
	Keys       map[string]any
	mu         sync.RWMutex
	Data       any
	HTMLRender render.HTMLRender
	FormMap    map[string]any
}

func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}
	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *Context) Get(key string) (value any, ok bool) {
	c.mu.RLock()
	value, ok = c.Keys[key]
	c.mu.RUnlock()
	return
}

// 进入对应路由的下一个方法
func (c *Context) Next() {
	c.Index++
	for c.Index < len(c.Handlers) && !c.IsAbort() {
		c.Handlers[c.Index](c)
		c.Index++
	}
}

func (c *Context) Abort() {
	c.Index = math.MaxInt
}

func (c *Context) IsAbort() bool {
	return c.Index == math.MaxInt || c.Index < -1
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) QueryDefault(key string, defaultVal string) string {
	query := c.R.URL.Query().Get(key)
	if stringUtils.IsBlank(query) {
		return defaultVal
	}
	return query
}

func (c *Context) QueryArray(key string) []string {
	query := c.R.URL.Query()
	return query[key]
}

func (c *Context) QueryMap() map[string][]string {
	return c.R.URL.Query()
}

func (c *Context) PostForm(key string) (any, error) {
	if res, ok := c.FormMap[key]; ok {
		return res, nil
	}
	if err := c.Bind(&c.FormMap); err != nil {
		return nil, err
	}
	return c.FormMap[key], nil
}

func (c *Context) parseMultipartForm() {
	if err := c.R.ParseMultipartForm(defaultMaxMemory); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			factory.Get(&logger.Log{}).Error(err.Error())
		}
	}
}

func (c *Context) FormFile(name string) *multipart.FileHeader {
	file, header, err := c.R.FormFile(name)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	return header
}

func (c *Context) FormFiles(name string) ([]*multipart.FileHeader, error) {
	multipartForm, err := c.MultipartForm()
	return multipartForm.File[name], err
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.R.ParseMultipartForm(defaultMaxMemory)
	return c.R.MultipartForm, err
}

func (c *Context) UploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func (c *Context) PathVariable(key string) any {
	if v, ok := c.Params[key]; ok {
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

func (c *Context) JSON(status int, data ...any) error {
	if len(data) == 1 {
		return c.render(status, &render.JSON{Data: data[0]})
	}
	return c.render(status, &render.JSON{Data: data})
}

func (c *Context) XML(status int, data any) error {
	return c.render(status, &render.XML{Data: data})
}

func (c *Context) Redirect(status int, url string) error {
	return c.render(status, &render.Redirect{Code: status, Request: c.R, Location: url})
}

func (c *Context) String(status int, format string, values ...any) error {
	return c.render(status, &render.String{Format: format, Data: values})
}

func (c *Context) Fail(code int, msg string) {
	c.String(code, msg)
}

func (c *Context) render(statusCode int, r render.Render) error {
	err := r.Render(c.W, statusCode)
	c.StatusCode = statusCode
	return err
}

func (c *Context) ParseTemplate(page string, data any) error {
	return c.render(http.StatusOK, &render.HTML{
		Data:       data,
		Name:       page,
		Template:   c.HTMLRender.Template,
		IsTemplate: true,
	})
}

func (c *Context) ToHTML(status int, html string) error {
	return c.render(status, &render.HTML{Data: html, IsTemplate: false})
}

func (c *Context) ParseTemplates(page string, data any, filenames ...string) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(page)
	t, err := t.ParseFiles(filenames...)
	if err != nil {
		return err
	}
	return t.Execute(c.W, data)
}

func (c *Context) FileDownload(filename string) {
	http.ServeFile(c.W, c.R, filename)
}

func (c *Context) FileDownloadWithFilename(filepath, filename string) {
	if stringUtils.IsASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.W, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: c.sameSite,
	})
}

func (c *Context) Bind(obj any) error {
	contentType := c.R.Header.Get("Content-Type")
	contentType = strings.Split(contentType, ";")[0]
	switch contentType {
	case binding.MIMEJSON:
		return c.BindJSON(obj)
	case binding.MIMEXML:
		return c.BindXML(obj)
	case binding.MIMEXML2:
		return c.BindXML2(obj)
	case binding.MIMEPLAIN:
		return c.BindPlain(obj)
	case binding.MIMEPOSTFORM:
		return c.BindFormPost(obj)
	case binding.MIMEMultipartPOSTForm:
		return c.BindMultipartPostForm(obj)
	case binding.MIMEMSGPACK2:
	case binding.MIMEMSGPACK:
		return c.BindMsgPack(obj)
	}
	return errors.New("unknown content-type : " + contentType)
}

func (c *Context) BindJSON(obj any) error {
	return binding.JSON.Bind(c.R, obj)
}

func (c *Context) BindXML(obj any) error {
	return binding.XML.Bind(c.R, obj)
}

func (c *Context) BindXML2(obj any) error {
	return binding.XML.Bind(c.R, obj)
}

func (c *Context) BindQuery(obj any) error {
	return binding.QUERY.Bind(c.R, obj)
}

func (c *Context) BindUri(obj any) error {
	uri := binding.URI
	uri.Params = c.Params
	return uri.Bind(c.R, obj)
}
func (c *Context) BindHeader(obj any) error {
	return binding.HEADER.Bind(c.R, obj)
}

func (c *Context) BindYAML(obj any) error {
	return binding.YAML.Bind(c.R, obj)
}

func (c *Context) BindPlain(obj any) error {
	return binding.PLAIN.Bind(c.R, obj)
}

func (c *Context) BindFormPost(obj any) error {
	return binding.FORM_POST.Bind(c.R, obj)
}

func (c *Context) BindMultipartPostForm(obj any) error {
	return binding.FORM_MULTIPART.Bind(c.R, obj)
}

func (c *Context) BindMsgPack(obj any) error {
	return binding.MSG_PACK.Bind(c.R, obj)
}
