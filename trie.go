package go_cloud

import (
	"context"
	"strings"
	"sync"
)

type Handler func(ctx *context.Context)

var GroupRouteIndex = 1

type Trie struct {
	CurNode   *Tree
	RootNode  *Tree
	Handlers  [][]Handler
	pool      sync.Pool
	PreIndex  int
	Index     int
	PreTrie   *Trie
	ParamsMap map[int][]string
}

type Tree struct {
	Next     map[string]*Tree
	Url      string
	End      bool
	Index    int
	PreIndex int
	PreTrie  *Trie
}

var RouteUrlParamsMap = make(map[int][]string)

var routeMap = make(map[string][]Handler)

var Root = &Tree{
	Next: make(map[string]*Tree),
}

var Estart = GetNewTrie()

func Init() *Trie {
	RouteInit(Root, nil, true)
	RouteDelete(Root, true)
	RouteUrlParamsMap = nil
	return Estart
}

func GetNewTrie() *Trie {
	trie := NewTrie()
	trie.pool.New = func() any {
		return &Context{}
	}
	return trie
}

func NewTrie() *Trie {
	return &Trie{
		Handlers:  make([][]Handler, 0),
		RootNode:  Root,
		CurNode:   Root,
		ParamsMap: make(map[int][]string),
	}
}

func (curNode *Tree) Insert(routeStringSlice []string, index int, handlerIndex int, e *Trie, routeIndex int) *Tree {
	curV := routeStringSlice[index]
	if []byte(routeStringSlice[index])[0] == ':' {
		curV = "/"
	}
	v, ok := curNode.Next[curV]
	if ok {
		curNode = v
	} else {
		newNode := &Tree{
			Next: make(map[string]*Tree),
			Url:  curV,
		}
		curNode.Next[newNode.Url] = newNode
		curNode = newNode
	}
	if curV == "/" {
		RouteUrlParamsMap[routeIndex] = append(RouteUrlParamsMap[routeIndex], string([]byte(routeStringSlice[index])[1:]))
	}
	if index+1 == len(routeStringSlice) {
		if handlerIndex == 0 && !curNode.End {
			curNode.End = false
			return curNode
		}
		curNode.PreIndex = routeIndex
		curNode.PreTrie = e
		curNode.End = true
		curNode.Index = handlerIndex
		return curNode
	}
	return curNode.Insert(routeStringSlice, index+1, handlerIndex, e, routeIndex)
}

func (e *Trie) Group(routeString string) *Trie {
	return e.GroupHandleIndex(routeString, 0)
}

func (e *Trie) GroupHandleIndex(routeString string, index int) *Trie {
	routeStringSlice := strings.Split(routeString, "/")
	start := 0
	end := len(routeStringSlice)
	if len(routeStringSlice[start]) == 0 {
		start = 1
	}
	if len(routeStringSlice[end-1]) == 0 {
		end = end - 1
	}
	newE := &Trie{
		CurNode:  e.CurNode,
		RootNode: Root,
		PreTrie:  e,
		Index:    GroupRouteIndex + 1,
	}
	GroupRouteIndex++
	curIndex := GroupRouteIndex
	newE.CurNode = newE.CurNode.Insert(routeStringSlice[start:end], 0, index, newE, curIndex)
	return newE
}

func (e *Trie) MatchRoute(routeString string) (bool, []Handler, map[string]any) {
	routeStringSlice := strings.Split(routeString, "/")
	e.CurNode = e.RootNode
	RouteParamMap := make(map[string]any)
	isMatch, handlerIndex, _ := e.CurNode.Match(routeStringSlice, 1, RouteParamMap, 0)
	if !isMatch {
		return false, nil, RouteParamMap
	}
	return isMatch, e.Handlers[handlerIndex-1], RouteParamMap
}

func (curNode *Tree) Match(routeStringSlice []string, index int, RouteParamMap map[string]any, urlIndex int) (bool, int, int) {
	if len(routeStringSlice) == index {
		if curNode.End {
			return true, curNode.Index, curNode.PreIndex
		}
		return false, 0, 0
	}
	tempNode1, tempNode2 := curNode, curNode
	v, ok := tempNode1.Next[routeStringSlice[index]]
	if ok {
		tempNode1 = v
		isMatch, handlerIndex, routeIndex := tempNode1.Match(routeStringSlice, index+1, RouteParamMap, urlIndex)
		if isMatch {
			return isMatch, handlerIndex, routeIndex
		}
	}
	v, ok = tempNode2.Next["/"]
	if ok {
		urlIndex++
		tempNode2 = v
		isMatch, handlerIndex, routeIndex := tempNode2.Match(routeStringSlice, index+1, RouteParamMap, urlIndex)
		if isMatch {
			RouteParamMap[Estart.ParamsMap[routeIndex][urlIndex-1]] = routeStringSlice[index]
			return isMatch, handlerIndex, routeIndex
		}
		urlIndex--
	}
	return false, 0, 0

}

func (e *Trie) AddRoute(routeName string, handler ...Handler) *Trie {
	Estart.Handlers = append(Estart.Handlers, handler)
	e = e.GroupHandleIndex(routeName, len(Estart.Handlers))
	return e
}

func RouteInit(root *Tree, routeSlice []string, mode bool) {
	if root == nil {
		return
	}
	if root.End {
		newHandlers, newPathParam := CallBackTree(root.PreTrie)
		Estart.Handlers[root.Index-1] = append(newHandlers, Estart.Handlers[root.Index-1]...)
		Estart.ParamsMap[root.PreIndex] = append(newPathParam, Estart.ParamsMap[root.PreIndex]...)
		if mode {
			routeMap["/"+strings.Join(routeSlice, "/")] = Estart.Handlers[root.Index-1]
		}
	}

	for k, v := range root.Next {
		if k != "/" {
			RouteInit(v, append(routeSlice, k), mode)
		} else {
			RouteInit(v, append(routeSlice, k), false)
		}
	}
}

func RouteDelete(root *Tree, mode bool) {
	if root.End {
		e := root.PreTrie
		if e != nil {
			if e.PreIndex > 0 {
				Estart.Handlers[e.PreIndex-1] = nil
			}
			e = e.PreTrie
		}
		if mode {
			root = nil
		}
		return
	}
	for k, v := range root.Next {
		if k != "/" {
			RouteDelete(v, mode)
		} else {
			RouteDelete(v, false)
		}
	}
}

func CallBackTree(e *Trie) ([]Handler, []string) {
	if e == nil {
		return nil, nil
	}
	var handlers []Handler
	var strArr []string
	if e.PreIndex > 0 {
		handlers = Estart.Handlers[e.PreIndex-1]
	}
	if e.Index > 0 {
		strArr = RouteUrlParamsMap[e.Index]
	}
	newHandlers, newPathParam := CallBackTree(e.PreTrie)
	handlers = append(handlers, newHandlers...)
	strArr = append(strArr, newPathParam...)
	return handlers, strArr
}
