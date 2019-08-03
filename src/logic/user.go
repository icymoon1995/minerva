package logic

import (
	"github.com/labstack/gommon/log"
	"minerva/src/common"
	logic "minerva/src/logic/common"
	"minerva/src/model"
)

/**
  相当于传统的UserService
     防止和 服务service 名字冲突
*/
type UserLogic struct {
}

/**
	校验用户名密码
    @param email string 邮箱
	@param password string 密码
	@return bool
*/
func (UserLogic) Verify(email string, password string) bool {
	// 应该连数据库去校验用户和密码
	//
	if email == "superhero" && password == "superhero,too" {
		return true
	}
	return false

}

/**
 *	登录
 */
func (UserLogic) Login() {
	// 设置cookie 、 利用 redis 存储缓存等
}

/**
 *	查询列表
 *	@params params map[string]interface{} 筛选的参数
 *	@return map[string]interface{}
 */
func (UserLogic) Index(params map[string]interface{}) map[string]interface{} {

	perPage, exist := params["per_page"].(int)
	if !exist {
		perPage = 10
	}
	currentPage, exist := params["page"].(int)
	if !exist {
		currentPage = 1
	}

	// 处理 params参数

	// 取总数量 - 分页使用
	user := &model.User{}
	total, _ := common.DB.Count(user)

	// 利用limit做分页处理
	offset := logic.Offset(currentPage, perPage)
	// 查询 (结合params加where条件)
	userList := make([]model.User, 0)
	error := common.DB.Limit(perPage, offset).Find(&userList)

	if error != nil {
		log.Print("logic.user#Index: ", error)
	}

	// 为了统一处理分页。。转换成interface{}
	data := make([]interface{}, perPage)
	for key, value := range userList {
		var temp interface{} = value
		data[key] = temp
	}

	return logic.Paginate(total, currentPage, perPage, data)

}

/**
 *	查询详情
 *	@param id int
 *	@return model.User
 */
func (UserLogic) Detail(id int) *model.User {
	user := &model.User{}
	_, error := common.DB.Id(id).Get(user)
	if error != nil {
		log.Print("logic.user#Detail error :", error)
	}
	return user
}
