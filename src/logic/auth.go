package logic

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"log"
	"minerva/src/common"
	"minerva/src/model"
)

var cookieSession = sessions.NewCookieStore([]byte("cookie_secret"))

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
	// 数据库去校验用户和密码
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

	// 通过hash加密作对比
	shaPassword := fmt.Sprintf("%X", sha1.Sum([]byte(password)))
	if shaPassword != auth.Password {
		return false, errors.New("password is not correct")
	}
	// setAuthInfo()
	return true, nil
}

/**
 *	登录
 */
func (AuthLogic) SetLoginInfo(ctx echo.Context, value string) {
	cookieSession.Options.HttpOnly = true
	// 设置cookie 、 利用 redis 存储缓存用户信息等
	session, _ := cookieSession.Get(ctx.Request(), "auth")
	session.Values["email"] = value
	_ = session.Save(ctx.Request(), ctx.Response())
}

func (AuthLogic) setAuthInfo() {

}
