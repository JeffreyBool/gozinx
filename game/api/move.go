/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 04:03
  * Software: GoLand
*/

package api

import (
	"github.com/JeffreyBool/gozinx/src/znet/router"
	"github.com/JeffreyBool/gozinx/src/ziface"
	pb "github.com/JeffreyBool/gozinx/game/proto"
	"github.com/gogo/protobuf/proto"
	"fmt"
	"github.com/JeffreyBool/gozinx/game/core"
)

/**
 * 玩家移动
 */

type MoveApi struct {
	router.BaseRouter
}

//路由之前的钩子方法
func (router *MoveApi) BeforeHandle(request ziface.IRequest) {

}

//路由处理业务方法
func (router *MoveApi) Handle(request ziface.IRequest) {
	//解析客户端传递过来的 proto 协议
	msg := &pb.Position{}
	if err := proto.Unmarshal(request.GetData(), msg); err != nil {
		fmt.Println("Move Position Unmarshal error: ", err)
		return
	}

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pid中获取
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		request.GetConnection().Stop()
		return
	}

	fmt.Printf("user pid = %d , move(%f,%f,%f,%f)", pid, msg.X, msg.Y, msg.Z, msg.V)

	//3. 根据pid得到player对象
	player := core.WorldManagerObj.GetPlayerByPid(pid.(uint32))

	//4. 让player对象发起移动位置信息广播
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}

//路由之后的钩子方法
func (router *MoveApi) AfterHandle(request ziface.IRequest) {

}
