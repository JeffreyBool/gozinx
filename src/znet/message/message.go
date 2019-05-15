/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 19:37
  * Software: GoLand
*/

package message

type Message struct {
	Id   uint32 //消息 id
	Size uint32 //消息长度
	Data []byte //消息内容
}

//获取消息 id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//获取消息长度
func (m *Message) GetMsgSize() uint32 {
	return m.Size
}

//获取消息内容
func (m *Message) GetMsg() []byte {
	return m.Data
}

//设置消息 id
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

//设置消息长度
func (m *Message) SetMsgSize(size uint32) {
	m.Size = size
}

//设置消息内容
func (m *Message) SetMsg(data []byte) {
	m.Data = data
}
