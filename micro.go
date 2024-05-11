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
	requestHandler *web.RequestHandler
}

func (e *Engine) Use(handler ...web.Handler) *Engine {
	e.requestHandler.UseGlobal(handler...)
	return e
}

func (e *Engine) Handle() *web.RequestHandler {
	return e.requestHandler
}

func New() *Engine {
	trie := web.NewTrie()
	return &Engine{
		trie: trie,
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
			fmt.Fprintln(w, r.RequestURI+" not found")
			return
		}
	}
	context = e.trie.GetEstart().Pool.Get().(*web.Context)
	context.R = r
	context.W = w
	context.Params = params
	context.Index = -1
	context.Handlers = handlers
	context.Next()
	w = context.W
	e.trie.GetEstart().Pool.Put(context)
}

func (e *Engine) Context() *web.Context {
	return e.trie.GetEstart().Pool.Get().(*web.Context)
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
