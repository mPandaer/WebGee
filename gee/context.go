package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//对于开发Web服务的人来说，提供原生的 ResponseWriter 和 Request粒度太细，每个Handler中都需要处理头部信息，比如Content-Type 以及响应码。
//所以我们需要封装一些东西，提供更加好用的API。提供常用的响应格式的API，而且除此之外，我们还需要提供关于本次请求的一些额外的信息，比如请求路径，
//而不是开发人员用底层Request来获取。

//基于此我们需要封装一些信息，所以我们的Context诞生了

var contentType = "Content-Type"
type JSON = map[string]interface{}

type Context struct {
	//提供底层的实例
	Writer http.ResponseWriter
	Req    *http.Request
	//请求信息
	Method string
	Path   string
	Params map[string]string

	//响应信息
	StatusCode int
}

func (c *Context) Param(key string) string {
	val,_ := c.Params[key]
	return val
}

// 向外暴露一个创建Context的函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
	}
}

//因为我们的Context更多的是提供好用的API，屏蔽用户对于ResponseWriter和Request的感知

// 提供获取请求参数的方法
func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

// 提供获取表单请求体的单值数据
func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

// 提供设置响应头的方法
func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Add(key, value)
}

// 提供设置响应状态码的方法
func (context *Context) SetStatus(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

//提供好用的响应API 处理方法

//响应HTML
func (context *Context) HTML(code int, html string) {
	context.SetHeader(contentType,"text/html")
	context.SetStatus(code)
	context.Writer.Write([]byte(html))
}

//响应JSON
func (context *Context) JSON(code int,data JSON) {
	context.SetStatus(code)
	context.SetHeader(contentType,"application/json")
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(context.Writer,err.Error(),500)
	}
}

//响应普通文本
func (context *Context) String(code int,format string, values ...interface{}) {
	context.SetStatus(code)
	context.SetHeader(contentType,"text/plain")
	context.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

//响应二进制数据
func (context *Context) Data(code int, data []byte) {
	context.SetStatus(code)
	context.Writer.Write(data)
}