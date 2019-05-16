/**
  * Author: JeffreyBool
  * Date: 2019/5/16
  * Time: 12:29
  * Software: GoLand
*/

package ziface

/**
 * 链接管理模块抽象定义
 */
type IConnManager interface {
	//得到当前链接总数
	Size() uint32

	//添加链接
	Add(connection IConnection)

	//删除链接
	Remove(connId uint32)

	//根据 ConnId 获取链接
	Get(connId uint32) (IConnection, error)

	//清除所有链接
	Clear()
}
