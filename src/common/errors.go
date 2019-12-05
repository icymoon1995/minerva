package common

import "errors"

var (
	NotFoundEmail      error = errors.New("账户不存在")
	AccountFrozen      error = errors.New("账户已冻结")
	PasswordNotCorrect error = errors.New("密码错误")
	ServerError        error = errors.New("服务器异常")
)

var (
	MqErrorFailConnect        = errors.New("Failed to connect to RabbitMq ")
	MqErrorChannelInitFail    = errors.New("Failed to initialize channel ")
	MqErrorQueueInitFail      = errors.New("Failed to initialize queue ")
	MqErrorExchangeInitFail   = errors.New("Failed to initialize exchange ")
	MqErrorQueueBindFail      = errors.New("Failed to bind queue ")
	MqErrorChannelConfirmFail = errors.New("Failed to open channel mode confirm  ")
	MqErrorSendMessageFail    = errors.New("Failed to send message after many tries ")
	MqErrorSendMessageTimeout = errors.New("send message timeout")
	MqMessageSuccess          = "send message success"
)
