package main

import (
	"net/http"

	"gee"
)

func main() {
	//创建我们的引擎对象
	engine := gee.New()

	//注册路由
	engine.GET("/hello", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>HELLO GEE</h1>")
	})

	g1 := engine.Group("/api")
	{
		g1.GET("/add", func(context *gee.Context) {
			context.JSON(http.StatusOK, gee.JSON{
				"username": "pandaer",
				"age":      19,
			})
		})
	}

	g2 := engine.Group("/admin")

	{
		g2.POST("/login", func(context *gee.Context) {
			username := context.PostForm("username")
			password := context.PostForm("password")
			if username == "pandaer" && password == "lwh" {
				context.JSON(http.StatusOK, gee.JSON{
					"res": "success",
				})
				return
			}
			context.JSON(http.StatusBadRequest, gee.JSON{
				"res": "fail",
				"msg": "密码或者用户名不正确",
			})
		})
	}

	//启动engine
	engine.Run(":9999")

}
