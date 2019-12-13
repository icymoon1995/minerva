package logic

import (
	"crypto/sha1"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"minerva/src/common"
	logic "minerva/src/logic/common"
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
	// 数据库去校验用户和密码
	auth := &model.Auth{}
	exist, error := common.DB.Where("email = ?", email).Get(auth)
	if error != nil {
		common.Logger.WithFields(logrus.Fields{
			"file":   "logic/auth.go",
			"method": "Verify",
			"type":   "verify error ",
		}).Errorln(error)
		// common.Logger.Errorln("logic.auth#Verify error :", error)
		return false, common.ServerError
	}

	if !exist {
		common.Logger.WithFields(logrus.Fields{
			"file":   "logic/auth.go",
			"method": "Verify",
			"type":   common.NotFoundEmail,
		}).Errorln("not found email (" + email + ")")
		// common.Logger.Errorln("logic.auth#Verify error : not found email (" + email + ")")
		return false, common.NotFoundEmail
	}

	// 通过hash加密作对比
	shaPassword := fmt.Sprintf("%X", sha1.Sum([]byte(password)))
	if shaPassword != auth.Password {
		common.Logger.WithFields(logrus.Fields{
			"file":   "logic/auth.go",
			"method": "Verify",
			"type":   "password error",
		}).Errorln("password is wrong")
		//common.Logger.Errorln("logic.auth#Verify error : password error (" + email + ")")
		return false, common.PasswordNotCorrect
	}

	// 账户激活验证
	if !auth.IsActive() {
		common.Logger.WithFields(logrus.Fields{
			"file":   "logic/auth.go",
			"method": "Verify",
			"type":   common.AccountFrozen,
		}).Errorln("account is frozen (" + email + ")")
		//	common.Logger.Errorln("logic.auth#Verify error : account is frozen : (" + email + ")")
		return false, common.AccountFrozen
	}

	// setAuthInfo()
	return true, nil
}

/**
 *	登录
 */
func (AuthLogic) SetLoginInfo(ctx echo.Context, value string) {
	// 保护cookie的安全性
	logic.CookieSession.Options.HttpOnly = true
	// 设置cookie 、 利用 redis 存储缓存用户信息等
	session, _ := logic.CookieSession.Get(ctx.Request(), "auth")
	session.Values["email"] = value
	_ = session.Save(ctx.Request(), ctx.Response())
}
