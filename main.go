package main

import (
	"fmt"
	"net/http"

	"gee"
)

func main() {
	//创建我们的引擎对象
	engine := gee.NewEngine()
	//注册路由
	engine.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello This is Gee\nlocation: %s\n", req.URL.Path)
	})

	engine.GET("/header", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	
	})

	//启动engine
	engine.Run(":9999")

}
