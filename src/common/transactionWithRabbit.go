package common

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"time"
)

/**
使用事务的小demo:

// 测试事务
//exchangeName := "transaction_exchange"
//routeKey := "key"
content := make(map[string]interface{})
content["name"] = "test_name"
content["email"] = "email@email.com"
message := common.Message{
	Id:       1,
	Action:   "add",
	Content:  content,
	Callback: "minerva/haha", // 消费成功后 调用的回调函数
}
// exchangeName,routeKey
err := common.TryMessageTransaction()
if err != nil {
	log.Println(err)
}
//  exchangeName, routeKey
err = common.MakeMessageTransaction(message)
if err != nil {
	log.Println(err)
}
*/

// 发送消息最大尝试次数
var MaxTries = 2

// 发送消息最大等待时间
var MaxTime = 700

// 消息动作 - try
var MessageActionPrepare = "prepare_ping"

// 专门处理事务的交换器
var TransactionExchangeName string

// 专门处理事务的路由key
var TransactionRouteKey string

/**
rabbitMq 事务初始化
*/
func transactionInit() {

	TransactionExchangeName = viper.GetString(rabbitPrefix + "transaction.exchangeName")
	TransactionRouteKey = viper.GetString(rabbitPrefix + "transaction.routeKey")
	var exchangeType string = viper.GetString(rabbitPrefix + "transaction.exchangeType")
	var queueName string = viper.GetString(rabbitPrefix + "transaction.queueName")

	// 创建exchange
	exchangeInit(TransactionExchangeName, exchangeType)

	// 创建队列
	queueInit(queueName, "", "")

	// 队列绑定
	queueBind(queueName, TransactionRouteKey, TransactionExchangeName)

}

// 默认无参数
/**
事务开启 - 尝试与队列通信
@return error
*/
func TryMessageTransaction() error {
	// 走默认的
	return TryMessageTransactionWithExchangeAndRoute(TransactionExchangeName, TransactionRouteKey)
}

/**
事务开启 - 尝试与队列通信
@param exchange string 交换器名
@param routeKey string 路由key
@return error
*/
func TryMessageTransactionWithExchangeAndRoute(exchange string, routeKey string) error {
	// send prepare
	content := make(map[string]interface{})
	prepareMessage := Message{
		Id:       -1,
		Action:   MessageActionPrepare,
		Content:  content,
		Callback: "", // 消费成功后 调用的回调函数
	}
	_, err := sendMessage(prepareMessage, exchange, routeKey)

	if err != nil {
		log.Println(prepareMessage.Action, err.Error())
		// panic(err)
	}
	return err

}

/**
发送事务处理的消息
@param trulyMessage Message 消息
@return error
*/
func MakeMessageTransaction(trulyMessage Message) error {
	return MakeMessageTransactionWithExchangeAndRoute(trulyMessage, TransactionExchangeName, TransactionRouteKey)
}

/**
发送事务处理的消息
@param trulyMessage Message 消息
@param exchange string 交换器名
@param routeKey string 路由key
@return error
*/
func MakeMessageTransactionWithExchangeAndRoute(trulyMessage Message, exchange string, routeKey string) error {
	_, err := sendMessage(trulyMessage, exchange, routeKey)

	if err != nil {
		log.Println(trulyMessage.Action, err.Error())
	}
	return err
}

/**
发送消息 (含重试处理)
@param message []byte 消息
@param exchange string 交换器
@param routeKey string 路由key
@return *amqp.Confimation, error
*/
func sendMessage(message Message, exchange string, routeKey string) (*amqp.Confirmation, error) {

	jsonMsg, jsonErr := json.Marshal(message)

	if jsonErr != nil {
		return nil, jsonErr
	}

	// 发送消息计数
	var count int = 0
	// 发送消息
	confirmation, err := mqSend(jsonMsg, exchange, routeKey)

	// 重试机制
	// 只有当消息超时再重新发送 , 其他错误直接返回
	for err != nil && err == MqErrorSendMessageTimeout {
		// 超过了最大重试次数
		if count >= MaxTries {
			return confirmation, MqErrorSendMessageFail
		}
		confirmation, err = mqSend(jsonMsg, exchange, routeKey)
		count++
	}

	return confirmation, err
}

/**
mq发送消息
@param message []byte 消息
@param exchange string 交换器
@param routeKey string 路由key
@return *amqp.Confimation, error
*/
func mqSend(message []byte, exchange string, routeKey string) (*amqp.Confirmation, error) {
	err := Channel.Publish(
		exchange,
		routeKey,
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
			// MessageId: "",
		})

	if err != nil {
		return nil, err
		// panic(err)
	}

	// 结果处理
	select {

	// 确认处理
	case messageConfirmation := <-ConfirmChan: //如果有数据，下面打印。但是有可能ch一直没数据
		log.Printf("%s: %s: %s: %d ", "common.transactionWithRabbit#mqSend", MqMessageSuccess, message, messageConfirmation.DeliveryTag)
		return &messageConfirmation, nil

	// 超时处理
	case <-time.After(time.Duration(MaxTime) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后MaxTime毫秒超时
		log.Printf("%s: %s: %s ", "common.transactionWithRabbit#mqSend", MqErrorSendMessageTimeout.Error(), message)
		return nil, MqErrorSendMessageTimeout
	}

}
