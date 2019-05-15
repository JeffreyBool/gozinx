/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 19:33
  * Software: GoLand
*/

package ziface

type IMessage interface {
	//获取消息 id
	GetMsgId() uint32

	//获取消息长度
	GetMsgSize() uint32

	//获取消息的内容
	GetMsg() []byte

	//设置消息的 id
	SetMsgId(uint32)

	//设置消息的长度
	SetMsgSize(uint32)

	//设置消息的内容
	SetMsg([]byte)
}
