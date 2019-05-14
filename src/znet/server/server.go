/**
  * Author: JeffreyBool
  * Date: 2019/5/14
  * Time: 22:23
  * Software: GoLand
*/

package server

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"net"
	"fmt"
	"github.com/JeffreyBool/gozinx/src/znet/connection"
	"github.com/pkg/errors"
)

/**
 服务器模块
 */

type Server struct {
	Config
}

//服务器配置
type Config struct {
	//服务器名称
	Name string
	//服务器绑定的 ip 版本
	IPVersion string
	//服务器监听的 IP
	IP string
	//服务器监听的端口
	Port int
}

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		return errors.Wrapf(err, "write back buf")
	}

	//回显业务
	return nil
}

//初始化
func NewServer(args ...Config) ziface.IServer {
	var config Config
	if len(args) > 0 {
		config = args[0]
	} else {
		config = Config{Name: "", IPVersion: "tcp4", IP: "0.0.0.0", Port: 8999}
	}

	return &Server{config}
}

func (s *Server) Start() error {
	fmt.Printf("[start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		//Todo 获取tcp的 addr
		address := fmt.Sprintf("[%s]:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, address)
		if err != nil {
			fmt.Printf("resolve tcp address err: %s\n", err)
			return
			//return err ors.Wrap(err, "resolve tcp address error")
		}

		//Todo 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen: %s, err: %s \n", s.IPVersion, err)
			return
			//return errors.Wrapf(err, "listen: %s", s.IPVersion)
		}

		//Todo 阻塞的等待客户端连接，处理客户端连接业务 （读写）
		fmt.Printf("start GoZinx server success, name: %s Listenning...\n", s.Name)
		var ConnID uint32 = 1
		for {
			//Todo 如果有客户端链接过来，会阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err: %s\n", err)
				continue
			}

			defer conn.Close()

			//Todo 已经与客户端建立链接
			c := connection.NewConnection(conn, ConnID, CallbackToClient)
			go c.Start()

			ConnID ++
		}
	}()

	return nil
}

//程序运行
func (s *Server) Serve() {
	s.Start()

	//Todo 做一些启动服务器之后的额外业务

	//阻塞，防止进程结束
	select {}
}

func (s *Server) Stop() error {
	//Todo
	return nil
}
