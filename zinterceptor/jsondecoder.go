package zinterceptor

import (
	"bytes"
	"github.com/zenus/zinx/ziface"
	"github.com/zenus/zinx/zlog"
	"sync"
)

const (
	SEP = '#'
	BEP = '*'
)

type JsonDecoder struct {
	lock sync.Mutex
	in   []byte
}

func NewJsonDecoder() ziface.IFrameDecoder {
	return new(JsonDecoder)
}

func (d *JsonDecoder) decode(buf []byte) []byte {
	in := bytes.NewBuffer(buf)
	b, err := in.ReadByte()
	if err != nil || b != SEP {
		return nil
	}
	buff := make([]byte, 0)
	for {
		b, err := in.ReadByte()
		if err != nil {
			return nil
		}
		if b == SEP {
			break
		}
		if b != BEP {
			buff = append(buff, b)
		}
	}
	return buff
}

func (d *JsonDecoder) Decode(buff []byte) [][]byte {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.in = append(d.in, buff...)
	resp := make([][]byte, 0)

	for {
		arr := d.decode(d.in)

		if arr != nil {
			//证明已经解析出一个完整包
			resp = append(resp, arr)
			_size := len(arr) + 4
			//_len := len(this.in)
			//fmt.Println(_len)
			if _size > 0 {
				zlog.Ins().DebugF("read before %s \n", string(d.in))
				d.in = d.in[_size:]
				zlog.Ins().DebugF("read after %s \n", string(d.in))
			}
		} else {
			return resp
		}
	}
}
