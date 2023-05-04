package main

import (
	"github.com/zenus/zinx/examples/zinx_async_op/router"
	"github.com/zenus/zinx/ziface"
	"github.com/zenus/zinx/zlog"
	"github.com/zenus/zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	zlog.Debug("zinx_async_op OnConnectionAdd ===>")
}

func OnConnectionLost(conn ziface.IConnection) {
	zlog.Debug("zinx_async_op OnConnectionLost ===>")
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouter(1, &router.LoginRouter{})

	s.Serve()
}
