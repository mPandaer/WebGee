package main

import (
	"net/http"

	"gee"
)

func main() {
	//创建我们的引擎对象
	engine := gee.NewEngine()
	//注册路由
	engine.GET("/",func(context *gee.Context) {
		context.HTML(http.StatusOK,"<h1>HELLO GEE</h1>")
	})

	engine.POST("/header", func(context *gee.Context) {
		context.JSON(http.StatusOK,gee.JSON{
			"username":context.PostForm("username"),
			"password":context.PostForm("password"),
		})
	})

	//启动engine
	engine.Run(":9999")

}
