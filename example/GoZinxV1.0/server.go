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

//HelloRouter test 自定义路由
type HelloRouter struct {
	router.BaseRouter
}

func (router *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Printf("recv from client MsgId: %d, MsgData: %s\n", request.GetMsgId(), request.GetData())
	//选读取客户端发送的数据
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("Hello Welcome To GoZinx\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//new server 服务
	s := server.NewServer()

	//注册连接的 hook 钩子函数
	s.SetOnConnStart(func(conn ziface.IConnection) {
		fmt.Println("====> Do Connection Begin is Called ....")
		if err := conn.SendMsg(1, []byte("Do Connection Begin Called \n")); err != nil {
			fmt.Println(err)
		}

		//设置两个链接属性，在连接创建之后
		fmt.Println("Set conn Name, Home done!")
		conn.SetProperty("Name", "JeffreyBool")
		conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")

		err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
		if err != nil {
			fmt.Println(err)
		}
	})

	//连接断开之前需要执行的函数
	s.SetOnConnStop(func(conn ziface.IConnection) {
		fmt.Printf("ConnId: [%d], ConnAddr: [%s] is Lost ...\n", conn.GetConnID(), conn.RemoteAddr().String())

		//在连接销毁之前，查询conn的Name，Home属性
		if name, err := conn.GetProperty("Name"); err == nil {
			fmt.Println("Conn Property Name = ", name)
		}

		if home, err := conn.GetProperty("Home"); err == nil {
			fmt.Println("Conn Property Home = ", home)
		}
		fmt.Println("====> Do Connection Lost is Called ....")
	})

	//给 server 添加一个自定义的 router
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	//启动 server
	s.Serve()
}
