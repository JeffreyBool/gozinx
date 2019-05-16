/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 01:17
  * Software: GoLand
*/

package main

import (
	"github.com/JeffreyBool/gozinx/src/znet/server"
	"github.com/JeffreyBool/gozinx/src/ziface"
	"github.com/JeffreyBool/gozinx/game/core"
	"github.com/JeffreyBool/gozinx/game/api"
	"fmt"
	"github.com/kr/pretty"
)

func main() {
	//创建 gozinx server
	serve := server.NewServer()

	//连接创建钩子函数
	serve.SetOnConnStart(func(conn ziface.IConnection) {
		//创建一个 player 对象
		player := core.NewPlayer(conn)

		//给客户端发送 msgId:1 的消息: 同步当前玩家的 Id 给客户端
		player.SyncPid()

		//给客户端发送 msgId:200 的消息 ： 同步当前 player 的初始位置给客户端
		player.BroadCastStartPosition()

		//将当前新上线的玩家添加到 WorldManager 中
		core.WorldManagerObj.AddPlayer(player)

		//将改连接绑定一个 pid 玩家的 Id属性
		conn.SetProperty("pid", player.Pid)

		//==============同步周边玩家上线信息，与现实周边玩家信息========
		player.SyncSurrounding()
		//=======================================================

		fmt.Printf("<===== Player Id: [%d] Is Arrived ========\n", player.Pid)
	})

	//连接销毁钩子函数
	serve.SetOnConnStop(func(conn ziface.IConnection) {
		//获取当前连接的Pid属性
		pid, _ := conn.GetProperty("pid")
		//根据pid得到player对象
		player := core.WorldManagerObj.GetPlayerByPid(pid.(uint32))
		pretty.Printf("%# v pid:[%d]\n", player, pid)
		player.LostConnection()
		fmt.Printf("<=====  Player Id: [%d]  Is Closed ========\n", pid)
	})

	//注册世界聊天广播路由
	serve.AddRouter(2, &api.WorldChatApi{}) //聊天

	//移动路由
	serve.AddRouter(3, &api.MoveApi{}) //移动

	//启动服务
	serve.Serve()
}
