package go_cloud

import (
	"errors"
	"fmt"
	"github.com/caixr9527/go-cloud/web"
	"log"
	"net/http"
)

type Engine struct {
	trie           *web.Trie
	middleware     []web.Handler
	requestHandler *web.RequestHandler
}

func (e *Engine) Use(handler ...web.Handler) *Engine {
	e.middleware = append(e.middleware, handler...)
	return e
}

func (e *Engine) Handle() *web.RequestHandler {
	return e.requestHandler
}

func New() *Engine {
	trie := web.NewTrie()
	return &Engine{
		trie:       trie,
		middleware: make([]web.Handler, 0),
		requestHandler: &web.RequestHandler{
			Trie: trie,
		},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	routeString := r.RequestURI + "/" + r.Method

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
			fmt.Fprintln(w, r.RequestURI+"not found")
			return
		}
	}
	context = e.trie.GetEstart().Pool.Get().(*web.Context)
	context.R = r
	context.W = w
	context.PathParams = params
	context.Index = -1
	context.Handlers = handlers
	e.Next()
	context.Next()
	w = context.W
	e.trie.GetEstart().Pool.Put(context)
}

func (e *Engine) Next() {
	for _, handler := range e.middleware {
		handler(e.trie.GetEstart().Pool.Get().(*web.Context))
	}
}

func (e *Engine) Run(addr string) {
	e.trie.Initialization()
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
