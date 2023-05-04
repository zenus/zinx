package c_router

import (
	"github.com/zenus/zinx/ziface"
	"github.com/zenus/zinx/zlog"
	"github.com/zenus/zinx/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call HelloZinxRouter Handle")
	zlog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

}
