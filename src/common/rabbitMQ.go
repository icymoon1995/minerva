package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"net"
	"strings"
	"time"
)

// 全局的rabbitMQ
var rabbit *amqp.Connection
var Channel *amqp.Channel
var Queue amqp.Queue

var rabbitPrefix string // redis配置前缀

// 发送者确认的队列
var ConfirmChan chan amqp.Confirmation

// 默认最大尝试次数
var MaxTries int = 2

// 发送消息最大等待时间
var MaxTime int = 700

// ping消息的类型
var MessageActionPrepare string = "prepare_ping"

// 死信队列
var deadExchange string = "dead_message_exchange"

// 死信路由key
var deadRouteKey string = "dead_key"

var (
	MqErrorFailConnect        = errors.New("Failed to connect to RabbitMq ")
	MqErrorChannelInitFail    = errors.New("fail to initialize channel ")
	MqErrorQueueInitFail      = errors.New("Failed to initialize queue ")
	MqErrorExchangeInitFail   = errors.New("Failed to initialize exchange ")
	MqErrorQueueBindFail      = errors.New("Failed to bind queue ")
	MqErrorChannelConfirmFail = errors.New("Failed to open channel mode confirm  ")
	MqErrorSendMessageFail    = errors.New("Failed to send message after many tries ")
	MqErrorSendMessageTimeout = errors.New("send message timeout")
	MqMessageSuccess          = "send message success"
)

// 消息
type Message struct {
	Id       int                    `json:"id"`
	Action   string                 `json:"action"`   // 消息动作
	Content  map[string]interface{} `json:"content"`  // 具体的消息内容
	Callback string                 `json:"callback"` // 消费成功后 调用的回调函数
}

// 初始化前缀
func rabbitPreInit() {
	rabbitPrefix = Enviorment + ".rabbit."
}

/**
初始化rabbit队列
*/
func rabbitInit() {

	// 初始化前缀
	rabbitPreInit()

	// rabbit conn连接初始化
	rabbitConn()

	// rabbit 组件初始化
	componentInit()

}

// 组件初始化
func componentInit() {

	// 创建通道
	channelInit()

	// 队列名
	var queueName string = viper.GetString(rabbitPrefix + "queue")
	// 创建队列
	queueInit(queueName, deadExchange, deadRouteKey)

	// 绑定的route key
	var routeKeys string = viper.GetString(rabbitPrefix + "routeKey")
	// 绑定的 交换机
	var exchangeName string = viper.GetString(rabbitPrefix + "exchange")
	// 绑定队列
	queueBind(Queue.Name, routeKeys, exchangeName)
}

/**
rabbitMq连接初始化
*/
func rabbitConn() {

	// 连接方式 默认为tcp
	var username string = viper.GetString(rabbitPrefix + "username")
	// 地址
	var host string = viper.GetString(rabbitPrefix + "host")
	// Auth密码
	var password string = viper.GetString(rabbitPrefix + "password")
	// 端口
	var port string = viper.GetString(rabbitPrefix + "port")

	// 生成rabbit的连接语句
	var rabbitConfig string = fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, host, port)

	var err error
	rabbit, err = amqp.DialConfig(rabbitConfig, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
		Heartbeat: time.Duration(2 * time.Second),
	})

	if err != nil {
		Logger.Fatalf("%s: %s: %s", "common.rabbitMq#rabbitConn", MqErrorFailConnect, err)
	}

	Logger.Println("rabbit mq connection success")
}

/**
创建通道
*/
func channelInit() {
	var err error
	Channel, err = rabbit.Channel()

	if err != nil {
		Logger.Fatalf("%s: %s: %s", "common.rabbitMq#channelInit", MqErrorChannelInitFail, err)
	}
	Logger.Println("rabbit mq channel success")

	// 开启发送者确认模式
	err = Channel.Confirm(false)
	ConfirmChan = make(chan amqp.Confirmation)
	Channel.NotifyPublish(ConfirmChan)

	if err != nil {
		Logger.Fatalf("%s: %s: %s", "common.rabbitMq#channelInit", MqErrorChannelConfirmFail, err)
	}
}

/**
	创建通道
 	@param name string 交换器名
 	@param exchangeType string 交换器类型
*/
func exchangeInit(name string, exchangeType string) {
	var err error
	// 创建exchange
	err = Channel.ExchangeDeclare(
		name,         // exchange_name
		exchangeType, // exchange_type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		Logger.Fatalf("%s: %s: %s", "common.rabbitMq#exchangeInit", MqErrorExchangeInitFail.Error(), err)
	}

}

/**
	创建队列
  	@param name	string	队列名
	@param deadExchange string 死信交换器
	@param deadRouteKey string 死信路由key
*/
func queueInit(name string, deadExchange string, deadRouteKey string) {
	var err error

	args := amqp.Table{}
	// 绑定该队列的 死信route和key
	args["x-dead-letter-exchange"] = deadExchange
	args["x-dead-letter-routing-key"] = deadRouteKey

	Queue, err = Channel.QueueDeclare(
		name,  // name
		true,  // durable -- 是否持久化
		false, // delete when unused
		false, // exclusive -- 是否独占
		false, // no-wait -- 阻塞消息
		args,  // arguments
	)

	if err != nil {
		Logger.Fatalf("%s: %s: %s ", "common.rabbitMq#queueInit", MqErrorQueueInitFail, err)
	}
}

/**
队列绑定
@param queueName string 队列名
@param routeKeys string 绑定的路由key,可以多个逗号分割
@param exchangeName string 绑定的exchange名称
*/
func queueBind(queueName string, routeKeys string, exchangeName string) {

	var err error

	var routeNames []string

	// 可能有多个route key
	routeNames = strings.Split(routeKeys, ",")

	// 目前只绑定了一个exchange
	// 在其他项目提前初始化exchange！
	for _, routeName := range routeNames {
		err = Channel.QueueBind(
			Queue.Name,   // queue name
			routeName,    // routing key
			exchangeName, // exchange
			false,
			nil,
		)
		if err != nil {
			Logger.Fatalf("%s: %s: %s ", "common.rabbitMq#queueBind", MqErrorQueueBindFail, err)
		}
	}

}

/**
发送消息 (含重试处理)
@param message &Message 消息
@param exchange string 交换器
@param routeKey string 路由key
@return *amqp.Confimation, error
*/
func SendMessage(message *Message, exchange string, routeKey string) (*amqp.Confirmation, error) {

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

		// todo  mq拒绝接受消息处理 (confirmation.Ack == false)
		if confirmation != nil && confirmation.Ack == false {

		}
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
		Logger.Printf("%s: %s: %s: %d ", "common.transactionWithRabbit#mqSend", MqMessageSuccess, message, messageConfirmation.DeliveryTag)
		return &messageConfirmation, nil

	// 超时处理
	case <-time.After(time.Duration(MaxTime) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后MaxTime毫秒超时
		Logger.Printf("%s: %s: %s ", "common.transactionWithRabbit#mqSend", MqErrorSendMessageTimeout.Error(), message)
		return nil, MqErrorSendMessageTimeout
	}

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
	prepareMessage := &Message{
		Id:       -1,
		Action:   MessageActionPrepare,
		Content:  content,
		Callback: "", // 消费成功后 调用的回调函数
	}
	_, err := SendMessage(prepareMessage, exchange, routeKey)

	if err != nil {
		Logger.Error(prepareMessage.Action, err.Error())
		// panic(err)
	}
	return err
}
