/**
  * Author: JeffreyBool
  * Date: 2019/5/17
  * Time: 00:11
  * Software: GoLand
*/

package core_test

import (
	"testing"
	"github.com/JeffreyBool/gozinx/game/core"
	"fmt"
	"runtime"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := core.NewAOIManager(0,250, 5, 0,250, 5)
	fmt.Println(aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiMgr := core.NewAOIManager(0,250, 5, 0,250, 5)

	for k, _ := range aoiMgr.Grids {
		//得到当前格子周边的九宫格
		grids := aoiMgr.GetSurroundGridsByGid(k)
		//得到九宫格所有的IDs
		fmt.Println("gid : ", k, " grids len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.Id)
		}
		fmt.Printf("grid ID: %d, surrounding grid IDs are %v\n", k, gIDs)
	}
}

func TestAOIManager_GetGidByPos(t *testing.T) {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	fmt.Println(fmt.Sprintf("%s", buf[:n]))
}
