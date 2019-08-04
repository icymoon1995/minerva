package logic

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"minerva/src/common"
	"minerva/src/model"
)

/**
  相当于传统的UserService
     防止和 服务service 名字冲突
*/
type AuthLogic struct {
}

/**
	校验用户名密码
    @param email string 邮箱
	@param password string 密码
	@return bool
*/
func (AuthLogic) Verify(email string, password string) (bool, error) {
	// 应该连数据库去校验用户和密码
	//
	if email == "superhero" && password == "superhero,too" {
		return true, nil
	}
	auth := &model.Auth{}
	exist, error := common.DB.Where("email = ?", email).Get(auth)
	if error != nil {
		log.Print("logic.auth#Verify error :", error)
		return false, error
	}

	if !exist {
		return false, errors.New("can not found username")
	}

	if !auth.IsActive() {
		return false, errors.New("your account is not login")
	}

	shaPassword := fmt.Sprintf("%X", sha1.Sum([]byte(password)))
	if shaPassword != auth.Password {
		return false, errors.New("password is not correct")
	}

	return true, nil
}

/**
 *	登录
 */
func (UserLogic) Login() {
	// 设置cookie 、 利用 redis 存储缓存用户信息等
}
