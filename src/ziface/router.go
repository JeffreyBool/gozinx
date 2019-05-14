/**
  * Author: JeffreyBool
  * Date: 2019/5/15
  * Time: 03:44
  * Software: GoLand
*/

package ziface

/**
  路由的抽象接口
  路由里的数据都是 IRequest
 **/
type IRouter interface {
	//路由之前的钩子方法
	BeforeHandle(IRequest)

	//路由处理业务方法
	Handle(IRequest)

	//路由之后的钩子方法
	AfterHandle(IRequest)
}
