package web

import (
	"net/http"
	"strings"
)

type RequestHandler struct {
	Trie    *Trie
	newTrie *Trie
}

func (rh *RequestHandler) Group(path string) *RequestHandler {
	rh.newTrie = rh.Trie.Group(path)
	return rh
}

func (rh *RequestHandler) UseGlobal(handlers ...Handler) {
	rh.Trie.Use(handlers...)
}

func (rh *RequestHandler) Use(handlers ...Handler) {
	rh.newTrie.Use(handlers...)
}

func (rh *RequestHandler) GET(path string, handlers ...Handler) {
	path = path + "/" + http.MethodGet
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) HEAD(path string, handlers ...Handler) {
	path = path + "/" + http.MethodHead
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) POST(path string, handlers ...Handler) {
	path = path + "/" + http.MethodPost
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) PUT(path string, handlers ...Handler) {
	path = path + "/" + http.MethodPut
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) PATCH(path string, handlers ...Handler) {
	path = path + "/" + http.MethodPatch
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) DELETE(path string, handlers ...Handler) {
	path = path + "/" + http.MethodDelete
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) CONNECT(path string, handlers ...Handler) {
	path = path + "/" + http.MethodConnect
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) OPTIONS(path string, handlers ...Handler) {
	path = path + "/" + http.MethodOptions
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) TRACE(path string, handlers ...Handler) {
	path = path + "/" + http.MethodTrace
	rh.check(path)
	rh.newTrie.AddRoute(path, handlers...)
}

func (rh *RequestHandler) check(path string) {
	index := strings.LastIndex(path, "/")
	realPath := path[0:index]
	method := path[index+1:]
	prePath := "/" + rh.newTrie.CurNode.CurPath
	if search, _, _ := rh.Trie.GetEstart().Search(prePath + path); search {
		panic("Duplicate routing [" + prePath + realPath + "], method [" + method + "]")
	}
}
