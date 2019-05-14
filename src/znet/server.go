/**
  * Author: JeffreyBool
  * Date: 2019/5/14
  * Time: 22:23
  * Software: GoLand
*/

package znet

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
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

func (Server) Start() error {
	panic("implement me")
}

func (Server) Serve() {
	panic("implement me")
}

func (Server) Stop() error {
	panic("implement me")
}
