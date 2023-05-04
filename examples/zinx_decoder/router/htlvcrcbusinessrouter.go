package router

import (
	"encoding/hex"
	"github.com/zenus/zinx/zdecoder"
	"github.com/zenus/zinx/ziface"
	"github.com/zenus/zinx/zlog"
	"github.com/zenus/zinx/znet"
)

type HtlvCrcBusinessRouter struct {
	znet.BaseRouter
}

func (this *HtlvCrcBusinessRouter) Handle(request ziface.IRequest) {

	//MsgID
	msgID := request.GetMessage().GetMsgID()
	zlog.Ins().DebugF("Call HtlvCrcBusinessRouter Handle %d %s\n", msgID, hex.EncodeToString(request.GetMessage().GetData()))

	resp := request.GetResponse()
	if resp == nil {
		return
	}

	tlvData := resp.(zdecoder.HtlvCrcDecoder)

	zlog.Ins().DebugF("do msgid=0x10 data business %+v\n", tlvData)
}
