/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 02:57
  * Software: GoLand
*/

package core

import (
	"sync"
)

/**
 * 当前游戏的世界总管理模块
 */
type WorldManager struct {
	//AoiManager 当前世界地图 AOI 的管理模块
	AoiManager *AOIManager

	//当前在线的 Players 集合
	Players map[uint32]*Player

	mutex *sync.RWMutex
}

var WorldManagerObj *WorldManager

//初始化方法
func init() {
	WorldManagerObj = &WorldManager{
		AoiManager: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players:    make(map[uint32]*Player),
		mutex:      new(sync.RWMutex),
	}
}

//提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	wm.Players[player.Pid] = player

	//将player 添加到AOI网络规划中
	wm.AoiManager.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

//从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid uint32) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	player, ok := wm.Players[pid]
	if ok {
		wm.AoiManager.RemoveFromGridByPos(int(pid), player.X, player.Z)
	}
	delete(wm.Players, pid)
}

//通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPid(pid uint32) *Player {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	return wm.Players[pid]
}

//获取所有玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range wm.Players {
		players = append(players, v)
	}

	//返回
	return players
}
