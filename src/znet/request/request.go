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
	conn ziface.IConnection

	//客户端请求的数据
	data []byte
}

//获取客户端链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//获取客户端数据
func (r *Request) GetData() []byte {
	return r.data
}
