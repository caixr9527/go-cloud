package go_cloud

import (
	"log"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type routerGroup struct {
	name           string
	handlerFuncMap map[string]HandlerFunc
}
type router struct {
	routerGroup []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	routerGroup := &routerGroup{
		name:           name,
		handlerFuncMap: make(map[string]HandlerFunc),
	}
	r.routerGroup = append(r.routerGroup, routerGroup)
	return routerGroup
}

func (r *routerGroup) Add(name string, handlerFunc HandlerFunc) {
	r.handlerFuncMap[name] = handlerFunc
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{},
	}
}

func (e *Engine) Run(addr string) {
	for _, group := range e.routerGroup {
		for key, value := range group.handlerFuncMap {
			http.HandleFunc(key, value)
		}
	}

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		// todo log
		log.Fatal(err)
	}
}
