package zpack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/zenus/zinx/zconf"
	"github.com/zenus/zinx/ziface"
)

// DataPackLtv
// Zinx早期使用的LTV 小端方式，兼容之前的应用
type DataPackLtv struct{}

// NewDataPack 封包拆包实例初始化方法
func NewDataPackLtv() ziface.IDataPack {
	return &DataPackLtv{}
}

// GetHeadLen 获取包头长度方法
func (dp *DataPackLtv) GetHeadLen() uint32 {
	//ID uint32(4字节) +  DataLen uint32(4字节)
	return defaultHeaderLen
}

// Pack 封包方法(压缩数据)
func (dp *DataPackLtv) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetRawData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法(解压数据)
func (dp *DataPackLtv) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if zconf.GlobalObject.MaxPacketSize > 0 && msg.GetDataLen() > zconf.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
