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
	"io"
)

/**
 链接模块
 */
type Connection struct {
	//当前链接的 socket tcp 套接字
	Conn *net.TCPConn

	//链接的 ID
	ConnID uint32

	//当前链接所绑定的处理业务方法回调函数
	Callback ziface.HandleFunc

	//当前链接的状态
	Close bool

	//告知当前链接已经退出 (close channel)
	Exit  chan bool
	mutex *sync.Mutex
}

//初始化链接
func NewConnection(conn *net.TCPConn, ConnID uint32, Callback ziface.HandleFunc) ziface.IConnection {
	return &Connection{
		Conn:     conn,
		ConnID:   ConnID,
		Callback: Callback,
		Exit:     make(chan bool, 1),
		mutex:    new(sync.Mutex),
	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn Start() ... ConnId = %d \n", c.ConnID)

	//启动从当前链接读取数据
	go c.startRead()

	//启动从当前链接写数据的业务
}

//链接读取
func (c *Connection) startRead() {
	fmt.Println("Conn Read Goroutine is running...")
	defer func() {
		fmt.Printf("Conn Stop ConnID = %d, Read is exit, remote addr is %s", c.ConnID, c.RemoteAddr().String())
		c.Stop()
	}()

	for {
		//读取客户端的数据到 buf 中，最大 512 字节
		buf := make([]byte, 512)
		read, err := c.Conn.Read(buf)
		if err == io.EOF {
			fmt.Printf("ConnID: %d exit\n", c.ConnID)
			return
		} else if err != nil {
			fmt.Printf("recv buf err: %s\n", err)
			continue
		}

		//调用客户端传递的回调函数
		if err = c.Callback(c.Conn, buf, read); err != nil {
			fmt.Println("ConnId", c.ConnID, "handle is error", err)
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

	c.Conn.Close()
	c.Close = true
	close(c.Exit)
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

func (c *Connection) Send([]byte) error {
	return nil
}
