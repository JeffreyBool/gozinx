/**
  * Author: JeffreyBool
  * Date: 2019/5/14
  * Time: 22:23
  * Software: GoLand
*/

package znet

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"fmt"
	"net"
)

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
		for {
			//Todo 如果有客户端链接过来，会阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err: %s\n", err)
				continue
			}

			defer conn.Close()

			//Todo 已经与客户端建立链接
			go func() {
				for {
					buf := make([]byte, 512)
					read, err := conn.Read(buf)
					if err != nil {
						fmt.Printf("recv buf err: %s \n", err)
						continue
					}

					//回写功能
					fmt.Printf("server name: %s client buf: %s, len: %d \n",s.Name,buf,read)
					if _, err = conn.Write(buf[:read]); err != nil {
						fmt.Printf("write back buf err: %s \n", err)
					}
				}
			}()
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
