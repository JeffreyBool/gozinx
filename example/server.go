/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 01:42
  * Software: GoLand
*/

package main

import (
	"github.com/JeffreyBool/gozinx/src/znet/server"
	"github.com/JeffreyBool/gozinx/src/znet/router"
	"github.com/go/src/fmt"
	"github.com/JeffreyBool/gozinx/src/ziface"
)

/*
 基于 `GoZinx` 框架来开发的，服务器端应用程序
**/

//ping test 自定义路由
type PingRouter struct {
	router.BaseRouter
}

func (router *PingRouter) BeforeHandle(request ziface.IRequest) {
	fmt.Println("Call Router BeforeHandle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...")); err != nil {
		fmt.Println("call back before ping error: ", err)
	}
}

func (router *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("handle ping...")); err != nil {
		fmt.Println("call back handle ping error: ", err)
	}
}

func (router *PingRouter) AfterHandle(request ziface.IRequest) {
	fmt.Println("Call Router AfterHandle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...")); err != nil {
		fmt.Println("call back after ping error: ", err)
	}
}

func main() {
	//启动 server 服务
	s := server.NewServer(server.Config{Name: "GoZinx V0.1", IPVersion: "tcp4", IP: "0.0.0.0", Port: 8999})
	s.Serve()
	s.AddRouter(&PingRouter{})
}
