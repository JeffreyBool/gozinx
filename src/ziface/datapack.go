/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 19:46
  * Software: GoLand
*/

package ziface

/**
 封包 、拆包 模块
 直接面向 tcp 连接中的数据流，用户处理 tcp 粘包的问题
**/

type IDataPack interface {
	//获取包头长度方法
	GetHeadSize() uint32

	//封包方法
	Pack(IMessage) ([]byte, error)

	//拆包方法
	Unpack([]byte) (IMessage, error)
}
