/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 01:49
  * Software: GoLand
*/

/*
 模拟客户端
**/

package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	//flag.String("server name","GoZinx V0.1",)

	fmt.Println("client start...")
	//Todo 直接链接服务器，得到一个 conn 链接
	conn, err := net.Dial("tcp", "[127.0.0.1]:8999")
	if err != nil {
		panic(err)
	}

	//Todo 释放链接句柄
	defer conn.Close()

	//Todo 链接调用 writ 写数据
	for {
		select {
		case <-time.After(time.Second):
			if _, err := conn.Write([]byte("hello word")); err != nil {
				fmt.Printf("write conn err: %s\n", err)
				return
			}

			buf := make([]byte, 512)
			if _, err := conn.Read(buf); err != nil {
				fmt.Printf("read buf err: %s\n", err)
				return
			}

			fmt.Printf("server call back: %s\n", buf)
		}
	}
}
