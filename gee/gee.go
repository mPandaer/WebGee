package gee

import (
	"fmt"
	"net/http"
)

// 我们应该能够拦截全部的HTTP请求，所以我们需要实现一个接口http.Handler. 又因为我们需要保存 Path-处理方法 的映射关系，所以我们可以定一个结构体
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

type engine struct {
	router map[string]HandlerFunc //这里暂时简单的处理一下，用一张HashMap保存映射关系，以后改进
}

// 向外提供一个创建Engine的方法
func NewEngine() *engine {
	return &engine{
		router: make(map[string]HandlerFunc),
	}
}

// 实现一个接口，从而拥有拦截所有HTTP请求的能力
// 当HTTP请求来了，会调用这个方法，我们这个时候就需要解析路由，并调用对应的处理方法
func (engine *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path //通过这种方式来唯一标识一个请求
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		//没有对应的处理方法。
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// 向外暴露一个注册路由的方法
func (engine *engine) AddRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// 封装注册路由的方法，简化常用注册逻辑。
func (engine *engine) GET(pattern string, handler HandlerFunc) {
	engine.AddRouter("GET", pattern, handler)
}

func (engine *engine) POST(pattern string, handler HandlerFunc) {
	engine.AddRouter("POST", pattern, handler)
}


//向外暴露一个开启HTTP服务的方法
func (engine *engine) Run(addr string) {
	http.ListenAndServe(addr,engine)
}