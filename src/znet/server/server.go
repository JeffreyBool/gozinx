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
	"github.com/JeffreyBool/gozinx/src/utils"
)

/**
 服务器模块
 **/

var zinx_logo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var top_line = `┌───────────────────────────────────────────────────┐`
var border_line = `│`
var bottom_line = `└───────────────────────────────────────────────────┘`

func init() {
	fmt.Println(top_line)
	fmt.Println(zinx_logo)
	fmt.Println(bottom_line)
}

type Server struct {
	Config

	//当前的 Server 添加一个 router, server 注册的链接对应的处理业务
	Router ziface.IRouter
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
		config = Config{Name: utils.GlobalObject.Name, IPVersion: "tcp4", IP: utils.GlobalObject.Host, Port: utils.GlobalObject.TcpPort}
	}

	return &Server{Config: config, Router: nil}
}

func (s *Server) Start() error {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	go func() {
		//获取tcp的 addr
		address := fmt.Sprintf("[%s]:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, address)
		if err != nil {
			fmt.Printf("resolve tcp address err: %s\n", err)
			return
			//return err ors.Wrap(err, "resolve tcp address error")
		}

		//监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen: %s, err: %s \n", s.IPVersion, err)
			return
			//return errors.Wrapf(err, "listen: %s", s.IPVersion)
		}

		//阻塞的等待客户端连接，处理客户端连接业务 （读写）
		fmt.Printf("start GoZinx server success, name: %s Listenning...\n", s.Name)
		var ConnID uint32 = 1
		for {
			//如果有客户端链接过来，会阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("Accept err: %s\n", err)
				continue
			}

			defer conn.Close()

			//已经与客户端建立链接
			c := connection.NewConnection(conn, ConnID, s.Router)
			go c.Start()

			ConnID ++
		}
	}()

	return nil
}

//程序运行
func (s *Server) Serve() {
	s.Start()

	//做一些启动服务器之后的额外业务

	//阻塞，防止进程结束
	select {}
}

//服务停止
func (s *Server) Stop() error {
	return nil
}

//服务添加路由
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}
