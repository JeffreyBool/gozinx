/**
  * Author: JeffreyBool
  * Date: 2019/5/14
  * Time: 22:23
  * Software: GoLand
*/

package ziface

//定义服务器抽象接口
type IServer interface {
	//启动服务器
	Start() error

	//运行服务器
	Serve()

	//停止服务器
	Stop() error

	//路由功能： 给当前的服务注册一个路由方法，供客户端链接处理方法
	AddRouter(uint32, IRouter) error

	//获取当前服务的连接管理器
	GetConnManager() IConnManager

	//注册 OnConnStart 钩子函数方法
	SetOnConnStart(ConnFunc)

	//注册 OnConnStop 钩子函数方法
	SetOnConnStop(ConnFunc)

	//调用 OnConnStart 钩子函数方法
	CallOnConnStart(conn IConnection)

	//调用 OnConnStop 钩子函数方法
	CallOnConnStop(conn IConnection)
}
