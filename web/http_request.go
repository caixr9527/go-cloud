package web

import "net/http"

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
	rh.newTrie.AddRoute(path+"/"+http.MethodGet, handlers...)
}

func (rh *RequestHandler) HEAD(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodHead, handlers...)
}

func (rh *RequestHandler) POST(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodPost, handlers...)
}

func (rh *RequestHandler) PUT(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodPut, handlers...)
}

func (rh *RequestHandler) PATCH(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodPatch, handlers...)
}

func (rh *RequestHandler) DELETE(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodDelete, handlers...)
}

func (rh *RequestHandler) CONNECT(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodConnect, handlers...)
}

func (rh *RequestHandler) OPTIONS(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodOptions, handlers...)
}

func (rh *RequestHandler) TRACE(path string, handlers ...Handler) {
	rh.newTrie.AddRoute(path+"/"+http.MethodTrace, handlers...)
}
