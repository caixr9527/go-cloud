package cloud

import (
	"errors"
	"fmt"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/internal/middleware"
	_ "github.com/caixr9527/go-cloud/log"
	"github.com/caixr9527/go-cloud/web"
	"github.com/caixr9527/go-cloud/web/render"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"sort"
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
	configuration := factory.Get(config.Configuration{})
	if configuration.Server.Https.Enable {
		e.runTLS(configuration)
	} else {
		e.run(configuration)
	}
}
func (e *Engine) run(configuration config.Configuration) {
	addr := fmt.Sprintf("%s%d", ":", configuration.Server.Port)
	e.trie.Initialization()
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	printLog(configuration, addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}

func initialization() {
	sort.Sort(component.Sort(component.Components))
	for index := range component.Components {
		component.Components[index].Create(component.SinglePool)
	}
}

func (e *Engine) runTLS(configuration config.Configuration) {
	logger := factory.Get(&zap.Logger{})
	addr := fmt.Sprintf("%s%d", ":", configuration.Server.Port)
	certFile := configuration.Server.Https.CertPath
	keyFile := configuration.Server.Https.KeyPath
	e.trie.Initialization()
	printLog(configuration, addr)
	logger.Info("load cert: " + certFile)
	logger.Info("load key: " + keyFile)
	err := http.ListenAndServeTLS(addr, certFile, keyFile, e)
	if err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}

func printLog(configuration config.Configuration, addr string) {
	fmt.Println("   _____  ____     _____ _      ____  _    _ _____  ")
	fmt.Println("  / ____|/ __ \\   / ____| |    / __ \\| |  | |  __ \\ ")
	fmt.Println(" | |  __| |  | | | |    | |   | |  | | |  | | |  | |")
	fmt.Println(" | | |_ | |  | | | |    | |   | |  | | |  | | |  | |")
	fmt.Println(" | |__| | |__| | | |____| |___| |__| | |__| | |__| |")
	fmt.Println("  \\_____|\\____/   \\_____|______\\____/ \\____/|_____/ " + Version)
	fmt.Println(" ::start on port" + addr)
	logger := factory.Get(&zap.Logger{})
	logger.Info("go-cloud start success, start on port" + addr)
	logger.Info("go-cloud env active: " + configuration.Cloud.Active)
	if configuration.Template.Path != "" {
		logger.Info("go-cloud load template: " + configuration.Template.Path)
	}
}

func (e *Engine) LoadTemplate(ops ...web.TemplateOps) {
	configuration := factory.Get(config.Configuration{})
	var funcMap template.FuncMap
	var pattern = configuration.Template.Path
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
