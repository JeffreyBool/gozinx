/**
  * Author: JeffreyBool
  * Date: 2019/5/16
  * Time: 16:19
  * Software: GoLand
*/

package core

import (
	"sync"
	"fmt"
)

/*
 *一个地图中的格子类
*/
type Grid struct {
	Id        int           //格子Id
	MinX      int           //格子左边界坐标
	MaxX      int           //格子右边界坐标
	MinY      int           //格子上边界坐标
	MaxY      int           //格子下边界坐标
	playerIds map[int]bool  //当前格子内的玩家或者物体成员ID
	mutex     *sync.RWMutex //读写锁
}

//初始化一个格子
func NewGrid(id, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		Id:        id,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIds: make(map[int]bool),
		mutex:     new(sync.RWMutex),
	}
}

//给格子添加一个玩家
func (g *Grid) Add(playerId int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.playerIds[playerId] = true
}

//从格子删除一个玩家
func (g *Grid) Remove(playerId int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.playerIds, playerId)
}

//得到当前格子中所有的玩家
func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for k, _ := range g.playerIds {
		playerIds = append(playerIds, k)
	}
	return
}

//打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v", g.Id, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIds)
}
