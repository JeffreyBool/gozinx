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
	"fmt"
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
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n")); err != nil {
		fmt.Println("call back before ping error: ", err)
	}
}

func (router *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("handle ping...\n")); err != nil {
		fmt.Println("call back handle ping error: ", err)
	}
}

func (router *PingRouter) AfterHandle(request ziface.IRequest) {
	fmt.Println("Call Router AfterHandle...")
	if _, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n")); err != nil {
		fmt.Println("call back after ping error: ", err)
	}
}

func main() {
	//new server 服务
	s := server.NewServer()

	//给 server 添加一个自定义的 router
	s.AddRouter(&PingRouter{})

	//启动 server
	s.Serve()
}
