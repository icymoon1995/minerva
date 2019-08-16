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
	// todo

	// 取总数量 - 分页使用
	user := &model.User{}
	total, _ := common.DB.Count(user)

	// 利用limit做分页处理
	offset := logic.Offset(currentPage, perPage)
	// 查询 (结合params加where条件)
	userList := make([]model.User, 0)
	// problem: 如果数据数量小于 perPage。。 则会返回多余很多null。。
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
