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

func (router *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Printf("recv from client MsgId: %d, MsgData: %s\n", request.GetMsgId(), request.GetData())
	//选读取客户端发送的数据
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("call server router ping...\n"))
	if err != nil {
		fmt.Println(err)
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
