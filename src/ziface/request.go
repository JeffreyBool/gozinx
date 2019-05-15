/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 03:38
  * Software: GoLand
*/

package ziface

/**
 将客户端请求的链接和数据包装到了一个 request 中
 */

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection

	//得到请求数据
	GetData() []byte

	//获取 msg Id
	GetMsgId() uint32
}
