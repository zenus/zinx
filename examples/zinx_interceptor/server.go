package main

import (
	"github.com/zenus/zinx/examples/zinx_interceptor/interceptors"
	"github.com/zenus/zinx/examples/zinx_interceptor/router"
	"github.com/zenus/zinx/znet"
)

func main() {
	server := znet.NewServer()

	server.AddRouter(1, &router.HelloRouter{})

	// Add Custom Interceptor
	server.AddInterceptor(&interceptors.MyInterceptor{})

	server.Serve()
}
