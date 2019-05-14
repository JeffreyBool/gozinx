/**
  * Author: JeffreyBool
  * Date: 2019/5/14
  * Time: 22:23
  * Software: GoLand
*/

package ziface

type IServer interface {
	//启动服务器
	Start() error
	//运行服务器
	Serve()
	//停止服务器
	Stop() error
}
