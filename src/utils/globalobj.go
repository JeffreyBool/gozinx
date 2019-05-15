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
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:           "GoZinx Server",
		Version:        "v0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalObject.reLoad()
}

//加载用户配置的文件
func (g *GlobalObj) reLoad() {
	data, err := ioutil.ReadFile("/Users/zhanggaoyuan/go/src/github.com/JeffreyBool/gozinx/src/conf/config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}
