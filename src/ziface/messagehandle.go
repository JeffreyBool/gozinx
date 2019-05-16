/**
  * Author: JeffreyBool
  * Date: 2019/5/16
  * Time: 01:56
  * Software: GoLand
*/

package ziface

/**
 * 消息处理接口抽象
 */
type IMessageHandle interface {
	//调度、执行对应的 router 消息处理方法
	DoMsgHandler(request IRequest) error

	//为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter) error

	//启动 worker 工作池
	StartWorkerPool()

	//将消息发送给消息队列处理
	SendMsgToTaskQueue(request IRequest)
}
