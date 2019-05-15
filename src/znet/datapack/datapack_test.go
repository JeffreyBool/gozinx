/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 20:40
  * Software: GoLand
*/

package datapack_test

import (
	"testing"
	"github.com/JeffreyBool/gozinx/src/znet/datapack"
	"net"
	"fmt"
	"io"
	"github.com/JeffreyBool/gozinx/src/znet/message"
)

func TestNewDataPack(t *testing.T) {
	dataPack := datapack.NewDataPack()
	t.Log(dataPack)
}

func TestDataPack_Pack(t *testing.T) {
	/**
	  模拟的服务器
	 */

	//创建 socket tcp
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error(err)
		return
	}

	go accept(listener)

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error(err)
		return
	}

	//创建一个封包对象 dp
	dp := datapack.NewDataPack()
	//封装一个msg1包
	msg := &message.Message{
		Id:   0,
		Size: 5,
		Data: []byte{'h', 'e', 'l', 'l', 'o'},
	}

	data, err := dp.Pack(msg)
	if err != nil {
		t.Error(err)
		return
	}

	msg = &message.Message{
		Id:   1,
		Size:7,
		Data:[]byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	data1, err := dp.Pack(msg)
	if err != nil {
		t.Error(err)
		return
	}

	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData := append(data, data1...)
	//向服务器端写数据
	conn.Write(sendData)

	//客户端阻塞
	select {}
}

//从客户端读取数据，拆包处理
func accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept error: ", err)
			return
		}

		defer conn.Close()

		go func(conn net.Conn) {
			//拆包
			dp := datapack.NewDataPack()
			for {
				buf := make([]byte, dp.GetHeadSize())
				_, err := io.ReadFull(conn, buf)
				if err != nil {
					fmt.Println("read head error: ", err)
					break
				}

				msg, err := dp.Unpack(buf)
				if err != nil {
					fmt.Println("server unpack error: ", err)
					break
				}

				if msg.GetMsgSize() > 0 {
					buf := make([]byte, msg.GetMsgSize())
					if _, err = io.ReadFull(conn, buf); err != nil {
						fmt.Println("server unpack data error: ", err)
						return
					}

					msg.SetMsg(buf)
					//完整的消息读取完毕，打印读取的消息
					fmt.Printf("recv msgId: %d, msgSize: %d, data: %s \n", msg.GetMsgId(), msg.GetMsgSize(), msg.GetMsg())
				}
			}
		}(conn)
	}
}
