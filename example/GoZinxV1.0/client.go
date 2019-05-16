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
	"github.com/JeffreyBool/gozinx/src/znet/message"
	"github.com/JeffreyBool/gozinx/src/znet/datapack"
	"io"
)

func main() {
	fmt.Println("client start...")
	//直接链接服务器，得到一个 conn 链接
	conn, err := net.Dial("tcp", "192.168.99.128:8999")
	if err != nil {
		panic(err)
	}

	//释放链接句柄
	defer conn.Close()

	//链接调用 writ 写数据
	for {
		select {
		case <-time.After(time.Second):

			//发送封包的 message 消息
			dp := datapack.NewDataPack()
			binaryMsg, err := dp.Pack(message.NewMessage(1, []byte("Zinx V1.0 Client Test Message")))
			if err != nil {
				fmt.Println("pack error: ", err)
				break
			}

			if _, err = conn.Write(binaryMsg); err != nil {
				fmt.Println("write error: ", err)
				break
			}

			//接受服务器返回的ping 消息拆包。
			headData := make([]byte, dp.GetHeadSize())
			if _, err = io.ReadFull(conn, headData); err != nil { //ReadFull 会把msg填充满为止
				fmt.Println("read head error: ", err)
				break
			}

			//将headData字节流 拆包到msg中
			msgHead, err := dp.Unpack(headData)
			if err != nil {
				fmt.Println("client unpack msg head error: ", err)
				break
			}

			if msgHead.GetMsgSize() > 0 {
				//再根据 data size 的长度将 data 读取出来
				buf := make([]byte, msgHead.GetMsgSize())
				if _, err := io.ReadFull(conn, buf); err != nil {
					fmt.Println("read msg data error: ", err)
					break
				}
				msgHead.SetMsg(buf)

				fmt.Printf("---> Recv Server MsgId: %d, MsgSize: %d, MsgData: %s \n", msgHead.GetMsgId(), msgHead.GetMsgSize(), msgHead.GetMsg())
			}
		}
	}
}
