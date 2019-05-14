/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 01:42
  * Software: GoLand
*/

package main

import (
	"github.com/JeffreyBool/gozinx/src/znet"
)

/*
 基于 `GoZinx` 框架来开发的，服务器端应用程序
**/

func main() {
	//启动 server 服务
	server := znet.NewServer(znet.Config{Name: "GoZinx V0.1", IPVersion: "tcp4", IP: "0.0.0.0", Port: 8999})
	server.Serve()
}
