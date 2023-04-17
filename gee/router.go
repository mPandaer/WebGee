package gee

import (
	"net/http"
	"strings"
)

//为了后续能支持动态路由，而且以后路由会越来越复杂，于是我们选择单独给路由的逻辑写一个文件，方便管理。

//路由的核心就是 提供对应请求的处理器方法，于是我们可以定义一个结构体来承载这份映射

// 处理器方法
type HandlerFunc func(context *Context)

type router struct {
	roots    map[string]*node // key对应的是请求方法，*node对应的是根节点。
	handlers map[string]HandlerFunc
}

//像外面暴露一个创建router的函数

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, part := range vs {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				return parts
			}
		}

	}
	return parts

}

// 注册路由的方法
func (router *router) AddRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	_, ok := router.roots[method]
	if !ok {
		router.roots[method] = &node{}
	}
	parts := parsePattern(pattern)
	router.roots[method].insert(pattern, parts, 0)
	router.handlers[key] = handler
}

func (router *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := router.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(pattern, searchParts, 0)

	if n == nil {
		return nil, nil
	}

	parts := parsePattern(n.pattern)
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]

		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}

	}

	return n, params

}

// 定义一个方法，用来处理调用对应的处理器方法
func (router *router) handle(context *Context) {
	//获取唯一的请求标识
	n, params := router.getRoute(context.Method, context.Path)
	if n != nil {
		context.Params = params
		key := context.Method + "-" + n.pattern
		router.handlers[key](context)
	} else {
		context.HTML(http.StatusNotFound, "<h1>404 NOT FOUND</h1>")
	}
}
