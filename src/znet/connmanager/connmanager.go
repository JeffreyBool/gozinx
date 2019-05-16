/**
  * Author: JeffreyBool
  * Date: 2019/5/16
  * Time: 12:34
  * Software: GoLand
*/

package connmanager

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"sync"
	"fmt"
	"github.com/pkg/errors"
)

/**
 * 连接管理模块
 */
type ConnManager struct {
	//连接 map 集合： key=> connId, value IConnection
	Connections map[uint32]ziface.IConnection

	//当前连接总数
	size uint32

	//读写锁
	lock *sync.RWMutex
}

func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
		lock:        new(sync.RWMutex),
	}
}

//获取当前连接总数
func (c *ConnManager) Size() uint32 {
	return c.size
}

//添加连接
func (c *ConnManager) Add(connection ziface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()

	//将 conn 加入到 ConnManager中
	c.Connections[connection.GetConnID()] = connection
	c.size ++
	fmt.Printf("add to connId: [%d] ConnManager successfully: conn num: [%d]", connection.GetConnID(), c.size)
}

//删除连接
func (c *ConnManager) Remove(connId uint32) {
	c.lock.Lock()
	defer c.lock.Unlock()

	//删除连接信息
	delete(c.Connections, connId)
	c.size --
	fmt.Printf("remove forom connId: [%d] ConnManager successfully: conn num: [%d]", connId, c.size)
}

//获取连接
func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.Connections[connId]
	if !ok {
		return conn, errors.Wrapf(errors.New("connection not found"), "connId: [%d", connId)
	}
	return conn, nil
}

//清空所有连接
func (c *ConnManager) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	//删除所有 conn并且停止工作
	for connId, conn := range c.Connections {
		//停止
		conn.Stop()
		//删除
		c.Remove(uint32(connId))
	}
}
