/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 03:15
  * Software: GoLand
*/

package api

import (
	"github.com/JeffreyBool/gozinx/src/znet/router"
	"github.com/JeffreyBool/gozinx/src/ziface"
	pd "github.com/JeffreyBool/gozinx/game/proto"
	"github.com/gogo/protobuf/proto"
	"fmt"
	"github.com/JeffreyBool/gozinx/game/core"
)

type WorldChatApi struct {
	router.BaseRouter
}

//路由处理业务方法
func (router *WorldChatApi) Handle(request ziface.IRequest) {
	//解析客户端传递过来的 proto 协议
	msg := &pd.Talk{}
	if err := proto.Unmarshal(request.GetData(), msg); err != nil{
		fmt.Println("Talk Unmarshal error: ",err)
		return
	}

	//当前的聊天数据是属于那个玩家发送的
	value, err := request.GetConnection().GetProperty("pid")
	if err != nil{
		fmt.Println("Get Property error: ",err)
		return
	}

	//根据 pid 得到对应的 player 对象
	player := core.WorldManagerObj.GetPlayerByPid(value.(uint32))

	//将这个消息广播给其他全部在线的玩家
	player.Talk(msg.Content)
}

