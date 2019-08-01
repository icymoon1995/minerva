package logic

import (
	"fmt"
	"minerva/src/common"
	"minerva/src/model"
)

/**
  相当于传统的UserService
     防止和 服务service 名字冲突
 */
type UserLogic struct {

}

func (u UserLogic) Detail (id int) *model.User {
	user := &model.User{}
	_, error := common.DB.Id(id).Get(user)
	if error != nil {
		fmt.Println("error : " , error)
	}
	return user
}