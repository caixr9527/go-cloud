package web

import "net/http"

type RequestHandler struct {
	Trie *Trie
}

func (rh *RequestHandler) Group(path string) *RequestHandler {
	rh.Trie.Group(path)
	return rh
}

func (rh *RequestHandler) Use(handlers ...Handler) *RequestHandler {
	rh.Trie.Use(handlers...)
	return rh
}

func (rh *RequestHandler) GET(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodGet, handlers...)
	return rh
}

func (rh *RequestHandler) HEAD(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodHead, handlers...)
	return rh
}

func (rh *RequestHandler) POST(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodPost, handlers...)
	return rh
}

func (rh *RequestHandler) PUT(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodPut, handlers...)
	return rh
}

func (rh *RequestHandler) PATCH(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodPatch, handlers...)
	return rh
}

func (rh *RequestHandler) DELETE(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodDelete, handlers...)
	return rh
}

func (rh *RequestHandler) CONNECT(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodConnect, handlers...)
	return rh
}

func (rh *RequestHandler) OPTIONS(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodOptions, handlers...)
	return rh
}

func (rh *RequestHandler) TRACE(path string, handlers ...Handler) *RequestHandler {
	rh.Trie.AddRoute(path+"/"+http.MethodTrace, handlers...)
	return rh
}
