/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 03:41
  * Software: GoLand
*/

package request

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
)

type Request struct {
	//客户端建立的链接
	Conn ziface.IConnection

	//客户端请求的数据
	Msg ziface.IMessage
}

//获取客户端链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

//获取客户端数据
func (r *Request) GetData() []byte {
	return r.Msg.GetMsg()
}

//获取消息 id
func (r *Request) GetMsgId() uint32 {
	return r.Msg.GetMsgId()
}
