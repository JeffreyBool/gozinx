/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 02:30
  * Software: GoLand
*/

package connection

import (
	"net"
	"github.com/JeffreyBool/gozinx/src/ziface"
	"fmt"
	"sync"
	"github.com/JeffreyBool/gozinx/src/znet/datapack"
	"io"
	"github.com/pkg/errors"
	"github.com/JeffreyBool/gozinx/src/znet/message"
	"github.com/JeffreyBool/gozinx/src/znet/request"
	"github.com/JeffreyBool/gozinx/src/utils"
)

/**
 链接模块
 */
type Connection struct {
	//当前 conn 属于哪个 server
	TcpServer ziface.IServer

	//当前链接的 socket tcp 套接字
	Conn *net.TCPConn

	//链接的 ID
	ConnID uint32

	//当前的 Server 的消息管理模块，用来绑定 msgId 和对应处理业务 API 关系
	MsgHandle ziface.IMessageHandle

	//当前链接的状态
	Close bool

	//用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//告知当前链接已经退出 (close channel)
	Exit chan bool

	//锁
	mutex *sync.Mutex
}

//初始化链接
func NewConnection(server ziface.IServer, conn *net.TCPConn, ConnID uint32, msgHandle ziface.IMessageHandle) ziface.IConnection {
	return &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    ConnID,
		MsgHandle: msgHandle,
		msgChan:   make(chan []byte),
		Exit:      make(chan bool, 1),
		mutex:     new(sync.Mutex),
	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn Start() ... ConnId = %d \n", c.ConnID)

	//添加一个连接
	c.TcpServer.GetConnManager().Add(c)

	//启动从当前链接读取数据
	go c.startReader()

	//启动从当前链接写数据的业务
	go c.StartWrite()
}

//读消息Goroutine，用于从客户端中读取数据
func (c *Connection) startReader() {
	fmt.Println("Conn Read Goroutine is running...")
	defer func() {
		fmt.Printf("Conn Stop ConnID = %d, Read is exit, remote addr is %s\n", c.ConnID, c.RemoteAddr().String())
		c.Stop()
	}()

	for {
		// 创建拆包解包的对象
		dp := datapack.NewDataPack()

		//读取客户端的 msg head 二进制流 8个字节
		buf := make([]byte, dp.GetHeadSize())
		if _, err := io.ReadFull(c.Conn, buf); err != nil {
			fmt.Println("read msg read error: ", err)
			break
		}

		//拆包，得到msgId和msg size放在 msg消息中
		message, err := dp.Unpack(buf)
		if err != nil {
			fmt.Println("conn unpack error: ", err)
			break
		}

		//根据 data size 再次读取data,放在 msg data 中。
		var data []byte
		if message.GetMsgSize() == 0 {
			fmt.Println("conn msg illegal data")
			break
		}

		data = make([]byte, message.GetMsgSize())
		if _, err = io.ReadFull(c.Conn, data); err != nil {
			fmt.Println("read msg data error: ", err)
			break
		}
		message.SetMsg(data)

		//得到当前 conn 数据的 Request 请求数据
		req := &request.Request{
			Conn: c,
			Msg:  message,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(req)
		} else {
			go func() {
				if err = c.MsgHandle.DoMsgHandler(req); err != nil {
					fmt.Println("do msg handle error: ", err)
					return
				}
			}()
		}
	}
}

//写消息Goroutine,专门发送给客户端消息
func (c *Connection) StartWrite() {
	fmt.Println("[Conn Write] Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "[Conn Write] exit")
	for {
		select {
		case msg, ok := <-c.msgChan:
			if !ok {
				return
			}

			if _, err := c.Conn.Write(msg); err != nil {
				fmt.Println("conn write error: ", err)
				return
			}
		case <-c.Exit:
			//代表 Reader 已经退出，此时 Write 也要退出
			return
		}
	}
}

//停止
func (c *Connection) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	fmt.Printf("Conn Stop() ... ConnId = %d\n", c.ConnID)

	if c.Close {
		return
	}

	c.Close = true
	c.Conn.Close()
	close(c.Exit)

	//删除当前连接
	c.TcpServer.GetConnManager().Remove(c.ConnID)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供 send msg方法，将我们要发送给客户端的数据进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.Close {
		return errors.New("connection closed when send msg")
	}

	//将 data 进行封包 ( size|id|data )
	dp := datapack.NewDataPack()
	binaryMsg, err := dp.Pack(message.NewMessage(msgId, data))
	if err != nil {
		return errors.Wrapf(err, "conn pack error msgId: %d, data: %s", msgId, data)
	}

	c.msgChan <- binaryMsg
	////将数据发送给客户端
	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	return errors.Wrap(err, "conn write error")
	//}

	return nil
}
