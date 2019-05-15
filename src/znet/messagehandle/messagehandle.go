/**
  * Author: JeffreyBool
  * Date: 2019/5/16
  * Time: 01:59
  * Software: GoLand
*/

package messagehandle

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"github.com/pkg/errors"
	"fmt"
)

/**
 * 消息处理模块的实现
 */
type MessageHandle struct {
	Routers map[uint32]ziface.IRouter
}

//初始化
func NewMessageHandle() ziface.IMessageHandle {
	return &MessageHandle{map[uint32]ziface.IRouter{}}
}

//调度 执行对应的 router 消息处理方法
func (m *MessageHandle) DoMsgHandler(request ziface.IRequest) error {
	route, ok := m.Routers[request.GetMsgId()]
	if ok {
		return errors.Wrapf(errors.New("do msg handle error"), "api MsgId: [%d] is not found need register")
	}

	route.BeforeHandle(request)
	route.Handle(request)
	route.AfterHandle(request)
	return nil
}

//消息添加路由关系
func (m *MessageHandle) AddRouter(msgId uint32, router ziface.IRouter) error {
	if _, ok := m.Routers[msgId]; ok {
		return errors.Wrapf(errors.New("repeat api error"), "MsgId: %d", msgId)
	}

	//添加msgId 与 api 的绑定关系
	m.Routers[msgId] = router
	fmt.Printf("add router api MsgId: [%d] success\n", msgId)
	return nil
}
