/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 18:50
  * Software: GoLand
*/

package utils

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"io/ioutil"
	"encoding/json"
)

/**
 存储一切有关 gozinx 框架的全局参数，供其他模块使用
 一些参数是可以通过 config.json 由用户进行配置
**/

type GlobalObj struct {
	//当前 GoZinx 全局的 Server 对象
	TcpServer ziface.IServer `json:"-"`

	//当前服务器主机监听的 ip
	Host string `json:"host"`

	//当前服务器主机监听的 端口号
	TcpPort int `json:"port"`

	//当前服务器的名称
	Name string `json:"name"`

	//gozinx 配置
	//当前gozinx的版本号
	Version string `json:"version"`

	//当前服务器主机允许的最大链接数
	MaxConn int `json:"max_conn"`

	//当前gozinx 框架数据包的最大值
	MaxPackageSize uint32 `json:"max_package_size"`

	//当前工作 worker 池的 goroutine 数量
	WorkerPoolSize uint32 `json:"worker_pool_size"`

	//GoZinx 框架允许用户最多开辟多少个 worker 数量
	MaxWorkerTaskSize uint32 `json:"max_worker_task_size"`
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:              "GoZinx Server",
		Version:           "v0.4",
		TcpPort:           8999,
		Host:              "0.0.0.0",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,
		MaxWorkerTaskSize: 1024,
	}

	GlobalObject.reLoad()
}

//加载用户配置的文件
func (g *GlobalObj) reLoad() {
	data, err := ioutil.ReadFile("conf/config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}
