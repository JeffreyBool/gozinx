/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 01:17
  * Software: GoLand
*/

package main

import (
	"github.com/JeffreyBool/gozinx/src/znet/server"
)

func main() {
	//创建 gozinx server
	serve := server.NewServer()

	//连接创建和销毁钩子函数

	//注册一些业务路由

	//启动服务
	serve.Serve()
}
