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
	"github.com/JeffreyBool/gozinx/src/utils"
)

/**
 * 消息处理模块的实现
 */
type MessageHandle struct {
	//负责 Worker 去任务的消息队列
	TaskQueue []chan ziface.IRequest

	//业务工作 Worker 池的数量
	WorkerPoolSize uint32

	//存放每个 msgId 对应的处理方法
	Routers map[uint32]ziface.IRouter
}

//初始化
func NewMessageHandle() ziface.IMessageHandle {
	return &MessageHandle{
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize),
		Routers:        map[uint32]ziface.IRouter{},
	}
}

//调度 执行对应的 router 消息处理方法
func (m *MessageHandle) DoMsgHandler(request ziface.IRequest) error {
	router, ok := m.Routers[request.GetMsgId()]
	if !ok {
		return errors.Wrapf(errors.New("do msg handle error"), "api MsgId: [%d] is not found need register", request.GetMsgId())
	}

	router.BeforeHandle(request)
	router.Handle(request)
	router.AfterHandle(request)
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
