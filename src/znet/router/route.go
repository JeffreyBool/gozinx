/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 03:44
  * Software: GoLand
*/

package router

import (
	"github.com/JeffreyBool/gozinx/src/ziface"
)

/**
  实现route时，先嵌入这个 BaseRouter 基类，然后根据需要对这个基类的方法进行重写就好了
 */
type BaseRouter struct {}

func (r *BaseRouter) BeforeHandle(ziface.IRequest) {}

func (r *BaseRouter) Handle(ziface.IRequest) {}

func (r *BaseRouter) AfterHandle(ziface.IRequest) {}

