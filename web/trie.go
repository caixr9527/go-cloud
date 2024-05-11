package web

import (
	"strings"
	"sync"
)

type Handler func(context *Context)

var GroupRouteIndex = 1

// 通过引擎来控制路由前缀树入口
type Trie struct {
	CurNode *TreeNode // 路由当前节点，比如在返回group函数后添加路由，每个group函数返回新的engine的当前节点
	// 的是前缀树对应的当前路由，group后添加的路由，从当前节点开始添加
	RootNode     *TreeNode   // 路由根节点
	HandlerSlice [][]Handler // 中间件，控制器方法数组
	Pool         sync.Pool
	PreIndex     int
	Index        int
	PreTrie      *Trie
	UrlParamsMap map[int][]string
}

// 前缀树节点，比如路由为/tee/api/:type/qq，那么路由会拆解成tee，api，:type，qq，四个节点,

type TreeNode struct {
	PathUrl  map[string]*TreeNode //当前路由对应下一个路由节点的url为key，下一个路由节点为value
	UrlValue string               // 当前路由url，例如api或qq
	// 当前路由为:type,:id等路径方式，存储参数，key为带有路由对应发放，
	// value为带有路径参数的url
	End      bool // 是否是路url由的重点
	Index    int  // 当前执行的路由方法索引
	PreIndex int  // 当前节点url对应的中间件

	PreTrie *Trie
}

// 为了支持路径参数，把带:的路由（例如:type)统一存储成/,并且用map存储执行函数索引和参数对应关系，存储到RouteUrlParamsMap
var routeUrlParamsMap = make(map[int][]string)
var routeMap = make(map[string][]Handler)

// 前缀树路由根节点
var root = &TreeNode{
	PathUrl: make(map[string]*TreeNode),
}

var eStart = GetNewTrie()

func Default() *Trie {
	return eStart
}

func (t *Trie) GetRouteMap() map[string][]Handler {
	return routeMap
}

func (t *Trie) GetEstart() *Trie {
	return eStart
}

// 插入一个前缀树路由
func (curNode *TreeNode) insert(routeStringSlice []string, index int, handlerIndex int, e *Trie, routeindex int) *TreeNode {
	curV := routeStringSlice[index]
	// 路径参数,保存为 /
	if []byte(routeStringSlice[index])[0] == ':' {
		curV = "/"
	}
	v, ok := curNode.PathUrl[curV]
	if ok {
		curNode = v
	} else {
		newNode := &TreeNode{
			PathUrl:  make(map[string]*TreeNode),
			UrlValue: curV,
		}

		// 上一个节点的map指向刚创建的节点
		curNode.PathUrl[newNode.UrlValue] = newNode
		curNode = newNode
	}
	// 如果是路由参数，对应方法索引指向该路径参数
	if curV == "/" {
		routeUrlParamsMap[routeindex] = append(routeUrlParamsMap[routeindex], string([]byte(routeStringSlice[index])[1:]))
	}

	// 路由插入前缀树完毕
	if index+1 == len(routeStringSlice) {
		// 如果是group，代表路由还没结束
		if handlerIndex == 0 && !curNode.End {
			curNode.End = false

			return curNode
		}
		curNode.PreIndex = routeindex
		// 如果是addRoute，代表路由到此结束，End设置为true，并保存路由索引
		curNode.PreTrie = e
		curNode.End = true
		curNode.Index = handlerIndex
		return curNode
	}

	return curNode.insert(routeStringSlice, index+1, handlerIndex, e, routeindex)

}

// 注册路由组到前缀树
func (t *Trie) Group(routeString string) *Trie {
	return t.groupHandleIndex(routeString, 0)
}

// 路由注册方法
func (t *Trie) groupHandleIndex(routeString string, handleIndex int) *Trie {
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
		CurNode:  t.CurNode,
		RootNode: root,
		PreTrie:  t,
		Index:    GroupRouteIndex + 1,
	}
	GroupRouteIndex++
	curIndex := GroupRouteIndex
	// 插入前缀路由，返回新的引擎
	newE.CurNode = newE.CurNode.insert(routeStringSlice[start:end], 0, handleIndex, newE, curIndex)

	return newE
}

// 匹配路由
func (t *Trie) Search(routeString string) (bool, []Handler, map[string]interface{}) {
	routeStringSlice := strings.Split(routeString, "/")
	t.CurNode = t.RootNode

	RouteParamMap := make(map[string]interface{})
	// isMatch为是否匹配，handerIndex为路由对应方法，handler为中间件
	isMatch, handlerIndex, _ := t.CurNode.Match(routeStringSlice, 1, RouteParamMap, 0)
	if !isMatch {
		return false, nil, RouteParamMap
	}

	return isMatch, t.HandlerSlice[handlerIndex-1], RouteParamMap
}

// 使用中间件
func (t *Trie) Use(handlers ...Handler) *Trie {

	if t.PreIndex != 0 {
		eStart.HandlerSlice[t.PreIndex-1] = append(eStart.HandlerSlice[t.PreIndex-1], handlers...)
	} else {
		eStart.HandlerSlice = append(eStart.HandlerSlice, handlers)
		t.PreIndex = len(eStart.HandlerSlice)
	}

	return t
}

// 匹配路由成功返回路径参数
func (curNode *TreeNode) Match(routeStringSlice []string, index int, RouteParamMap map[string]interface{}, urlindex int) (bool, int, int) {

	if len(routeStringSlice) == index {
		if curNode.End {
			return true, curNode.Index, curNode.PreIndex
		}
		return false, 0, 0
	}
	// 匹配不带路径参数的路由
	tempNode1, tempNode2 := curNode, curNode
	v, ok := tempNode1.PathUrl[routeStringSlice[index]]
	if ok {
		tempNode1 = v
		isMatch, handlerIndex, routeindex := tempNode1.Match(routeStringSlice, index+1, RouteParamMap, urlindex)
		if isMatch {
			return isMatch, handlerIndex, routeindex
		}
	}
	// 匹配带路径参数的路由
	v, ok = tempNode2.PathUrl["/"]
	if ok {
		urlindex++
		tempNode2 = v
		isMatch, handlerIndex, routeindex := tempNode2.Match(routeStringSlice, index+1, RouteParamMap, urlindex)
		if isMatch {

			RouteParamMap[eStart.UrlParamsMap[routeindex][urlindex-1]] = routeStringSlice[index]
			return isMatch, handlerIndex, routeindex
		}
		urlindex--
	}
	return false, 0, 0
}

// 增加路由
func (t *Trie) AddRoute(routeName string, handler ...Handler) *Trie {

	eStart.HandlerSlice = append(eStart.HandlerSlice, handler)
	t = t.groupHandleIndex(routeName, len(eStart.HandlerSlice))
	return t
}

func NewTrie() *Trie {
	return &Trie{
		HandlerSlice: make([][]Handler, 0),
		RootNode:     root,
		CurNode:      root,
		UrlParamsMap: make(map[int][]string),
	}

}

func GetNewTrie() *Trie {
	e := NewTrie()
	e.Pool.New = func() any {
		return &Context{}
	}
	return e
}

func (t *Trie) Initialization() {
	t.routeInit(root, nil, true)
	t.routeDelete(root, true)
	routeUrlParamsMap = nil
}

func (t *Trie) routeInit(root *TreeNode, routeSlice []string, mode bool) {
	if root == nil {
		return
	}
	if root.End {
		newHandlers, newPathParam := CallBackTreeNode(root.PreTrie)
		eStart.HandlerSlice[root.Index-1] = append(newHandlers, eStart.HandlerSlice[root.Index-1]...)

		eStart.UrlParamsMap[root.PreIndex] = append(newPathParam, eStart.UrlParamsMap[root.PreIndex]...)

		if mode {
			routeMap["/"+strings.Join(routeSlice, "/")] = eStart.HandlerSlice[root.Index-1]
		}

	}

	for k, v := range root.PathUrl {
		if k != "/" {
			t.routeInit(v, append(routeSlice, k), mode)
		} else {
			t.routeInit(v, append(routeSlice, k), false)
		}

	}
}

func (t *Trie) routeDelete(root *TreeNode, mode bool) {

	if root.End {
		e := root.PreTrie
		if e != nil {
			if e.PreIndex > 0 {
				eStart.HandlerSlice[e.PreIndex-1] = nil
			}
			e = e.PreTrie
		}
		if mode {
			root = nil
		}
		return
	}

	for k, v := range root.PathUrl {
		if k != "/" {
			t.routeDelete(v, mode)
		} else {
			t.routeDelete(v, false)

		}

	}

}

func CallBackTreeNode(e *Trie) ([]Handler, []string) {
	if e == nil {
		return nil, nil
	}
	var handlers []Handler
	var strArr []string
	if e.PreIndex > 0 {
		handlers = eStart.HandlerSlice[e.PreIndex-1]

	}

	if e.Index > 0 {
		strArr = routeUrlParamsMap[e.Index]
	}
	newHandlers, newPathParam := CallBackTreeNode(e.PreTrie)
	handlers = append(newHandlers, handlers...)
	strArr = append(newPathParam, strArr...)
	return handlers, strArr
}

func (t *Trie) GET(routeName string, handlers ...Handler) {
	t.AddRoute(routeName+"/GET", handlers...)
}

func (t *Trie) POST(routeName string, handlers ...Handler) {
	t.AddRoute(routeName+"/POST", handlers...)
}

func (t *Trie) DELETE(routeName string, handlers ...Handler) {
	t.AddRoute(routeName+"/DELETE", handlers...)
}
