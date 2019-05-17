/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 20:03
  * Software: GoLand
*/

package datapack

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"bytes"
	"encoding/binary"
	"github.com/JeffreyBool/gozinx/src/znet/message"
	"github.com/JeffreyBool/gozinx/src/utils"
	"github.com/pkg/errors"
)

/**
 * 封包、拆包具体模块
**/

const size uint32 = 8

type DataPack struct {
}

//拆包、封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的长度
func (data *DataPack) GetHeadSize() uint32 {
	// data size(4 byte) + id （4byte）
	return size
}

/**
 * 封包
 * data size|msgId|data
**/
func (data *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放 bytes 字节的缓冲
	buf := new(bytes.Buffer)

	//将 data size 写进 buf 中
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgSize()); err != nil {
		return nil, err
	}

	//将 msgId 写进 buf 中
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将 data 写进 buf 中
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsg()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/**
 * 拆包
 * 将包的 head 消息读出来，之后再根据 head 消息里的长度，再进行一次性读取
**/
func (data *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的 ioReader
	buf := bytes.NewReader(binaryData)

	//直接压 head 消息，得到的 data size 和 msgId
	msg := &message.Message{}

	//读取 data size
	if err := binary.Read(buf, binary.LittleEndian, &msg.Size); err != nil {
		return msg, err
	}

	//读取 msgId
	if err := binary.Read(buf, binary.LittleEndian, &msg.Id); err != nil {
		return msg, err
	}

	//判断当前包的长度是否超过了我们设置允许最大值
	if utils.GlobalObject.MaxPackageSize > 0 && msg.Size > utils.GlobalObject.MaxPackageSize {
		return msg, errors.New("too large msg data recv")
	}

	return msg, nil
}
