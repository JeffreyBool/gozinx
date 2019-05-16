/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 01:39
  * Software: GoLand
*/

package core

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
	"sync/atomic"
	"math/rand"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	pb "github.com/JeffreyBool/gozinx/game/proto"
)

var PidGen uint32 = 0 //用来生产玩家的 Id计数器

//玩家对象
type Player struct {
	Pid  uint32             //玩家 Id
	Conn ziface.IConnection //当前玩家的连接(用于和客户端的连接)
	X    float32            //平面的 X 坐标
	Y    float32            //高度
	Z    float32            //平面 Y 坐标
	V    float32            //旋转的 0-360 角度
}

//创建玩家方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成玩家的 Id
	atomic.AddUint32(&PidGen, 1)

	//创建玩家对象
	return &Player{
		Pid:  PidGen,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴偏移若干坐标
		Y:    0,                            //高度为0
		Z:    float32(134 + rand.Intn(17)), //随机在134坐标点 基于Y轴偏移若干坐标
		V:    0,                            //角度为0，尚未实现
	}
}

//告知客户端pid,同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组件 msgId: 1 的 proto 数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(1, data)
}

//广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//组件msgId: 200 的 proto 数据
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, data)
}

//玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	//1. 组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //TP 1 代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//2. 得到当前世界所有的在线玩家
	players := WorldManagerObj.GetAllPlayers()

	//3. 向所有的玩家发送MsgId:200消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

//给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
func (p *Player) SyncSurrounding() {
	//根据自己的位置，获取周围九宫格内的玩家pid
	pids := WorldManagerObj.AoiManager.GetPidsByPos(p.X, p.Z)
	//2 根据pid得到所有玩家对象
	players := make([]*Player, 0, len(pids))
	//3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
	for _, pid := range pids {
		players = append(players, WorldManagerObj.GetPlayerByPid(uint32(pid)))
	}

	//3.1 组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
	for _, player := range players {
		player.SendMsg(200, msg)
	}

	//4 让周围九宫格内的玩家出现在自己的视野中
	//4.1 制作Message SyncPlayers 数据
	playersData := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersData = append(playersData, p)
	}

	//4.2 封装SyncPlayer protobuf数据
	SyncPlayersMsg := &pb.SyncPlayers{
		Ps: playersData[:],
	}

	//4.3 给当前玩家发送需要显示周围的全部玩家数据
	p.SendMsg(202, SyncPlayersMsg)
}

//广播玩家位置移动
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	//更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组装protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4, //4 - 移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取当前玩家周边全部玩家
	players := p.GetSurroundingPlayers()
	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

//获得当前玩家的AOI周边玩家信息
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI区域的所有pid
	pids := WorldManagerObj.AoiManager.GetPidsByPos(p.X, p.Z)

	//将所有pid对应的Player放到Player切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldManagerObj.GetPlayerByPid(uint32(pid)))
	}

	return players
}

/*
	发送消息给客户端，
	主要是将pb的protobuf数据序列化之后发送
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) error {
	//将proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "player to sendMsg marshal error")
	}

	//判断客户端是否已经离线
	if p.Conn == nil {
		return errors.New("connection in player is nil")
	}

	//调用Zinx框架的SendMsg发包
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		return errors.New("player sendMsg error")
	}
	return nil
}
