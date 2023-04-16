package gee

import "net/http"

//为了后续能支持动态路由，而且以后路由会越来越复杂，于是我们选择单独给路由的逻辑写一个文件，方便管理。

//路由的核心就是 提供对应请求的处理器方法，于是我们可以定义一个结构体来承载这份映射

//处理器方法
type HandlerFunc func(context *Context)

type router struct {
	handlers map[string]HandlerFunc
}
//像外面暴露一个创建router的函数

func newRouter() *router{
	return &router{handlers: make(map[string]HandlerFunc)}
}


//注册路由的方法
func (router *router) AddRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handler
}

//定义一个方法，用来处理调用对应的处理器方法

func (router *router) handle(context *Context) {
	//获取唯一的请求标识
	key := context.Method + "-" + context.Path
	if handler,ok := router.handlers[key]; ok {
		handler(context)
	}else {
		context.HTML(http.StatusNotFound,"<h1>404 NOT FOUND</h1>")
	}
}

