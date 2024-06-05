package cloud

import (
	"errors"
	"fmt"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/internal/middleware"
	logger "github.com/caixr9527/go-cloud/log"
	"github.com/caixr9527/go-cloud/web"
	"github.com/caixr9527/go-cloud/web/render"
	"html/template"
	"log"
	"net/http"
)

type Engine struct {
	trie           *web.Trie
	requestHandler *web.RequestHandler
	ops            *web.Options
}

func (e *Engine) Use(handler ...web.Handler) *Engine {
	e.requestHandler.UseGlobal(handler...)
	return e
}

func (e *Engine) Handle() *web.RequestHandler {
	return e.requestHandler
}

func Default() *Engine {
	initialization()
	trie := web.NewTrie()
	return &Engine{
		trie: trie,
		requestHandler: &web.RequestHandler{
			Trie: trie,
		},
	}
}

func New(options *web.Options) *Engine {
	engine := Default()
	engine.Use(middleware.Recovery, middleware.Logging)
	engine.ops = options
	engine.LoadTemplate()
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	routeString := r.URL.Path + "/" + r.Method

	var context *web.Context
	var handlers []web.Handler
	var params map[string]any
	var isMatch bool
	if v, ok := e.trie.GetRouteMap()[routeString]; ok {
		handlers, params = v, make(map[string]any)
	} else {
		isMatch, handlers, params = e.trie.GetEstart().Search(routeString)
		if !isMatch {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, r.RequestURI+" not found")
			return
		}
	}
	context = e.trie.GetEstart().Pool.Get().(*web.Context)
	e.initContext(w, r, context, params, handlers)
	context.Next()
	w = context.W
	e.trie.GetEstart().Pool.Put(context)
}

func (e *Engine) initContext(w http.ResponseWriter, r *http.Request, context *web.Context, params map[string]any, handlers []web.Handler) {
	context.R = r
	context.W = w
	context.Params = params
	context.Index = -1
	context.Handlers = handlers
	context.HTMLRender = e.ops.HTMLRender
	context.FormMap = make(map[string]any)
}

func (e *Engine) Context() *web.Context {
	return e.trie.GetEstart().Pool.Get().(*web.Context)
}

func (e *Engine) Run() {
	if config.Cfg.Server.Https.Enable {
		e.runTLS()
	} else {
		e.run()
	}
}
func (e *Engine) run() {
	addr := fmt.Sprintf("%s%d", ":", config.Cfg.Server.Port)
	e.trie.Initialization()
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	printLog(addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}

func initialization() {
	config.Init()
	logger.Init()
}

func (e *Engine) runTLS() {
	addr := fmt.Sprintf("%s%d", ":", config.Cfg.Server.Port)
	certFile := config.Cfg.Server.Https.CertPath
	keyFile := config.Cfg.Server.Https.KeyPath
	e.trie.Initialization()
	printLog(addr)
	logger.Log.Info("load cert: " + certFile)
	logger.Log.Info("load key: " + keyFile)
	err := http.ListenAndServeTLS(addr, certFile, keyFile, e)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}

func printLog(addr string) {
	fmt.Println("   _____  ____     _____ _      ____  _    _ _____  ")
	fmt.Println("  / ____|/ __ \\   / ____| |    / __ \\| |  | |  __ \\ ")
	fmt.Println(" | |  __| |  | | | |    | |   | |  | | |  | | |  | |")
	fmt.Println(" | | |_ | |  | | | |    | |   | |  | | |  | | |  | |")
	fmt.Println(" | |__| | |__| | | |____| |___| |__| | |__| | |__| |")
	fmt.Println("  \\_____|\\____/   \\_____|______\\____/ \\____/|_____/ " + Version)
	fmt.Println(" ::start on port" + addr)
	logger.Log.Info("go-cloud start success, start on port" + addr)
	logger.Log.Info("go-cloud env active: " + config.Cfg.Cloud.Active)
	if config.Cfg.Template.Path != "" {
		logger.Log.Info("go-cloud load template: " + config.Cfg.Template.Path)
	}
}

func (e *Engine) LoadTemplate(ops ...web.TemplateOps) {
	var funcMap template.FuncMap
	var pattern = config.Cfg.Template.Path
	if len(ops) == 0 {
		funcMap = e.ops.FuncMap
		if e.ops.TemplatePattern != "" {
			pattern = e.ops.TemplatePattern
		}
	} else {
		funcMap = ops[0].FuncMap
		if ops[0].TemplatePattern != "" {
			pattern = ops[0].TemplatePattern
		}
	}
	t := template.Must(template.New("").Funcs(funcMap).ParseGlob(pattern))
	e.ops.HTMLRender = render.HTMLRender{Template: t}
}
