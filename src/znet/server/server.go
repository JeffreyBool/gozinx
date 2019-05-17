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
	"github.com/JeffreyBool/gozinx/src/znet/messagehandle"
	"github.com/JeffreyBool/gozinx/src/znet/connmanager"
)

/**
 服务器模块
 **/

var zinx_logo = `
   ██████╗  ██████╗ ███████╗██╗███╗   ██╗██╗  ██╗
 ██╔════╝ ██╔═══██╗╚══███╔╝██║████╗  ██║╚██╗██╔╝
 ██║  ███╗██║   ██║  ███╔╝ ██║██╔██╗ ██║ ╚███╔╝
 ██║   ██║██║   ██║ ███╔╝  ██║██║╚██╗██║ ██╔██╗
 ╚██████╔╝╚██████╔╝███████╗██║██║ ╚████║██╔╝ ██╗
  ╚═════╝  ╚═════╝ ╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝ `

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

	//当前的 Server 的消息管理模块，用来绑定 msgId 和对应处理业务 API 关系
	MsgHandle ziface.IMessageHandle

	//该 server 的连接管理器
	ConnManager ziface.IConnManager

	//创建连接之后自动调用 Hook 函数
	OnConnStart ziface.ConnFunc

	//销毁连接之前自动调用 Hook 函数
	OnConnStop ziface.ConnFunc
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

	return &Server{
		Config:      config,
		MsgHandle:   messagehandle.NewMessageHandle(),
		ConnManager: connmanager.NewConnManager(),
	}
}

func (s *Server) Start() error {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	go func() {
		//开启消息队列以及 worker 工作池
		s.MsgHandle.StartWorkerPool()

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

			//判断设置连接的最大值
			if int(s.ConnManager.Size()) >= utils.GlobalObject.MaxConn {
				fmt.Printf("Too Many Connections MaxConn: [%d] \n", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//已经与客户端建立链接
			c := connection.NewConnection(s, conn, ConnID, s.MsgHandle)
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
	//将一些服务器的资源、状态或者一些已经开辟的连接信息进行停止或者回收
	fmt.Printf("[Start] GoZinx Server Name: [%s] \n", s.Name)
	s.ConnManager.Clear()
	return nil
}

//服务添加路由
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) error {
	return s.MsgHandle.AddRouter(msgId, router)
}

//获取当前服务的连接管理器
func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

//注册 OnConnStart 钩子函数方法
func (s *Server) SetOnConnStart(hookFunc ziface.ConnFunc) {
	s.OnConnStart = hookFunc
}

//注册 OnConnStop 钩子函数方法
func (s *Server) SetOnConnStop(hookFunc ziface.ConnFunc) {
	s.OnConnStop = hookFunc
}

//调用 OnConnStart 钩子函数方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Printf("---> Cll OnConnStart() ....\n")
		s.OnConnStart(conn)
	}
}

//调用 OnConnStop 钩子函数方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Printf("---> Cll CallOnConnStop() ....\n")
		s.OnConnStop(conn)
	}
}
