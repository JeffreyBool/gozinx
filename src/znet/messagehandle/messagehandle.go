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
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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
	fmt.Printf("add router api MsgId: [%d] success \n", msgId)
	return nil
}

//启动 worker 工作池,
func (m *MessageHandle) StartWorkerPool() {
	for i := 0; i < int(utils.GlobalObject.WorkerPoolSize); i ++ {
		//当前的 worker 对应的 channel 消息队列开辟空间， 第 1 个就用 1 个 channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)

		//启动当前的 worker， 阻塞等待消息从 channel 传递过来
		go m.startWorker(i, m.TaskQueue[i])
	}
}

//启动 worker
func (m *MessageHandle) startWorker(workerId int, queue chan ziface.IRequest) {
	fmt.Printf("[Worker] Id: %d is started...\n", workerId)
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的 request,执行当前的 request 所绑定的业务
		case request := <-queue:
			m.DoMsgHandler(request)
		}
	}
}

//将消息交给 taskQueue, 由 worker 进行处理
func (m *MessageHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的 worker
	//根据客户端建立的 connId来进行分配
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Printf("Add ConnId: %d Request MsgId: %d, To WorkerId: %d \n", request.GetConnection().GetConnID(), request.GetMsgId(), workerId)
	//将消息发送给对应的 worker 的 taskQueue 即可
	m.TaskQueue[workerId] <- request
}
