package common

import "errors"

var (
	NotFoundEmail      error = errors.New("账户不存在")
	AccountFrozen      error = errors.New("账户已冻结")
	PasswordNotCorrect error = errors.New("密码错误")
	ServerError        error = errors.New("服务器异常")
)
