package gee

import (
	"log"
	"net/http"
)

// 路由分组控制
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine //保证在增加分组时，全局只有一个engine

	//新的一个分组，前缀支持嵌套
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix, //嵌套前缀
		middlewares: []HandlerFunc{},
		parent:      group,
		engine:      engine,
	}
	//创建了新的分组于是需要更新分组的集合。
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.AddRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(comp string, handler HandlerFunc) {
	group.addRoute("GET", comp, handler)
}

func (group *RouterGroup) POST(comp string, handler HandlerFunc) {
	group.addRoute("POST", comp, handler)
}

// 我们应该能够拦截全部的HTTP请求，所以我们需要实现一个接口http.Handler. 又因为我们需要保存 Path-处理方法 的映射关系，所以我们可以定一个结构体
// type HandlerFunc func(w http.ResponseWriter, req *http.Request)

type engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 向外提供一个创建Engine的方法
func New() *engine {
	engine := &engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine} //为什么这里没有给前缀，而是零值 “” 因为注册路由的时候都是以/开始的，所以似乎默认就实现了/的分组
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 实现一个接口，从而拥有拦截所有HTTP请求的能力
// 当HTTP请求来了，会调用这个方法，我们这个时候就需要解析路由，并调用对应的处理方法
func (engine *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := newContext(w, req)
	engine.router.handle(context)
}

// 向外暴露一个开启HTTP服务的方法
func (engine *engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
