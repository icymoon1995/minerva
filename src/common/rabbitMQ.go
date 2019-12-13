package common

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"net"
	"strings"
	"time"
)

var RabbitMq = &rabbit{}

// ping消息的类型
var MessageActionPrepare string = "prepare_ping"

// 配置
type rabbit struct {
	connection *amqp.Connection
	Channel    *amqp.Channel
	exchange   struct {
		name  string
		eType string
	}
	/**
	 多个exchange和多个队列的配置 怎么表现 todo
	exchange []struct {
		name string
		eType string
	}
	queueConfig []struct {
		queue amqp.Queue // 队列
		routeKeys []string // 绑定的路由key
		exchangeName string // 绑定的交换
	}
	*/
	QueueConfig struct {
		Queue        amqp.Queue // 队列
		routeKeys    []string   // 绑定的路由key
		exchangeName string     // 绑定的交换
	}
	confirmChan  chan amqp.Confirmation
	maxTries     int // 默认最大尝试次数
	maxTime      int // 发送消息最大等待时间
	deadExchange struct {
		name     string
		routeKey string
	}
	config struct { // 配置
		prefix string
		// 连接方式 默认为tcp
		username string
		// 地址
		host string
		// Auth密码
		password string
		// 端口
		port string
	}
}

// 消息
type Message struct {
	Id       int                    `json:"id"`
	Action   string                 `json:"action"`   // 消息动作
	Content  map[string]interface{} `json:"content"`  // 具体的消息内容
	Callback string                 `json:"callback"` // 消费成功后 调用的回调函数
}

func newRabbit() {
	RabbitMq.configInit()
	RabbitMq.rabbitConn()
	RabbitMq.channelInit()
	RabbitMq.queueInit()
	RabbitMq.queueBind()
	RabbitMq.exchangeInit()
}

/**
  配置初始化
*/
func (r *rabbit) configInit() {
	//  初始化前缀
	r.config.prefix = Global.Environment + ".rabbit."
	// 连接方式 默认为tcp
	r.config.username = viper.GetString(r.config.prefix + "username")
	// 地址
	r.config.host = viper.GetString(r.config.prefix + "host")
	// Auth密码
	r.config.password = viper.GetString(r.config.prefix + "password")
	// 端口
	r.config.port = viper.GetString(r.config.prefix + "port")

	// 死信路由key
	r.deadExchange.routeKey = "dead_key"
	// 死信队列
	r.deadExchange.name = "dead_message_exchange"

}

/**
  获取连接url
*/
func (r *rabbit) connectionUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s",
		r.config.username,
		r.config.password,
		r.config.host,
		r.config.port,
	)
}

/**
rabbitMq连接初始化
*/
func (r *rabbit) rabbitConn() {
	r.configInit()

	// 生成rabbit的连接语句
	var connection string = r.connectionUrl()

	var err error
	r.connection, err = amqp.DialConfig(connection, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
		Heartbeat: time.Duration(2 * time.Second),
	})

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "rabbitConn",
			"error":  MqErrorFailConnect,
		}).Fatalln(err)
		// Logger.Fatalf("%s: %s: %s", "common.rabbitMq#rabbitConn", MqErrorFailConnect, err)
	}

	Logger.WithFields(logrus.Fields{
		"file":   "common/rabbitMQ.go",
		"method": "rabbitConn",
		"logger": "rabbit mq",
	}).Println("connection success")
}

/**
rabbitMq channel
*/
func (r *rabbit) channelInit() {
	var err error
	r.Channel, err = r.connection.Channel()

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "channelInit",
			"error":  MqErrorChannelInitFail,
		}).Fatalln(err)
		//Fatalf("%s: %s: %s", "common.rabbitMq#channelInit", MqErrorChannelInitFail, err)
	}
	Logger.WithFields(logrus.Fields{
		"file":   "common/rabbitMQ.go",
		"method": "channelInit",
		"logger": "rabbit mq channel",
	}).Println("connection success")

	// 开启发送者确认模式
	err = r.Channel.Confirm(false)
	r.confirmChan = make(chan amqp.Confirmation)
	r.Channel.NotifyPublish(r.confirmChan)

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "channelInit",
			"error":  MqErrorChannelConfirmFail,
		}).Fatalln(err)
		//Fatalf("%s: %s: %s", "common.rabbitMq#channelInit", MqErrorChannelConfirmFail, err)
	}
}

/**
创建队列
*/
func (r *rabbit) queueInit() {
	var err error

	args := amqp.Table{}
	// 绑定该队列的 死信route和key
	args["x-dead-letter-exchange"] = r.deadExchange.name
	args["x-dead-letter-routing-key"] = r.deadExchange.routeKey

	// 队列名
	var queueName string = viper.GetString(r.config.prefix + "queue")

	// 创建队列
	r.QueueConfig.Queue, err = r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable -- 是否持久化
		false,     // delete when unused
		false,     // exclusive -- 是否独占
		false,     // no-wait -- 阻塞消息
		args,      // arguments
	)

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "queueInit",
			"error":  MqErrorQueueInitFail,
		}).Fatalln(err)
		//Fatalf("%s: %s: %s ", "common.rabbitMq#queueInit", MqErrorQueueInitFail, err)
	}
}

/**
队列绑定
*/
func (r *rabbit) queueBind() {

	var err error
	// 绑定的route key
	var routeKeys string = viper.GetString(r.config.prefix + "routeKey")
	// 绑定的 交换机名
	var exchangeName string = viper.GetString(r.config.prefix + "exchange")

	var routeNames []string
	// 可能有多个route key
	routeNames = strings.Split(routeKeys, ",")

	r.QueueConfig.routeKeys = routeNames
	r.QueueConfig.exchangeName = exchangeName

	// 目前只绑定了一个exchange
	// 在其他项目提前初始化exchange！
	for _, routeName := range routeNames {
		err = r.Channel.QueueBind(
			r.QueueConfig.Queue.Name,   // queue name
			routeName,                  // routing key
			r.QueueConfig.exchangeName, // exchange
			false,
			nil,
		)
		if err != nil {
			Logger.WithFields(logrus.Fields{
				"file":   "common/rabbitMQ.go",
				"method": "queueBind",
				"error":  MqErrorQueueBindFail,
			}).Fatalln(err)
			//Fatalf("%s: %s: %s ", "common.rabbitMq#queueBind", MqErrorQueueBindFail, err)
		}
	}

}

/**
创建通道
*/
func (r *rabbit) exchangeInit() {
	// 绑定的route key
	r.exchange.name = viper.GetString(r.config.prefix + "exchange")
	r.exchange.eType = "direct"

	var err error
	// 创建exchange
	err = r.Channel.ExchangeDeclare(
		r.exchange.name,  // exchange_name
		r.exchange.eType, // exchange_type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "queueBind",
			"error":  MqErrorExchangeInitFail,
		}).Fatalln(err)
		//Fatalf("%s: %s: %s", "common.rabbitMq#exchangeInit", MqErrorExchangeInitFail.Error(), err)
	}

}

/**
发送消息 (含重试处理)
@param message &Message 消息
@param exchange string 交换器
@param routeKey string 路由key
@return *amqp.Confimation, error
*/
func (r *rabbit) SendMessage(message *Message, exchange string, routeKey string) (*amqp.Confirmation, error) {

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
		if count >= r.maxTries {
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
func (r *rabbit) mqSend(message []byte, exchange string, routeKey string) (*amqp.Confirmation, error) {
	err := r.Channel.Publish(
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
	case messageConfirmation := <-r.confirmChan: //如果有数据，下面打印。但是有可能ch一直没数据
		//	Logger.Printf("%s: %s: %s: %d ", "common.transactionWithRabbit#mqSend", MqMessageSuccess, message, messageConfirmation.DeliveryTag)
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "mqSend",
			"type":   MqMessageSuccess,
		}).Println(message, messageConfirmation.DeliveryTag)
		return &messageConfirmation, nil

	// 超时处理
	case <-time.After(time.Duration(r.maxTime) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后MaxTime毫秒超时
		//	Logger.Printf("%s: %s: %s ", "common.transactionWithRabbit#mqSend", MqErrorSendMessageTimeout.Error(), message)
		Logger.WithFields(logrus.Fields{
			"file":   "common/rabbitMQ.go",
			"method": "mqSend",
			"type":   MqErrorSendMessageTimeout.Error(),
		}).Println(message)
		return nil, MqErrorSendMessageTimeout
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
	return RabbitMq.SendMessage(message, exchange, routeKey)
}

/**
mq发送消息
@param message []byte 消息
@param exchange string 交换器
@param routeKey string 路由key
@return *amqp.Confimation, error
*/
func mqSend(message []byte, exchange string, routeKey string) (*amqp.Confirmation, error) {
	return RabbitMq.mqSend(message, exchange, routeKey)
}

/**
事务开启 - 尝试与队列通信
@param exchange string 交换器名
@param routeKey string 路由key
@return error
*/
func (r *rabbit) TryMessageTransactionWithExchangeAndRoute(exchange string, routeKey string) error {
	// send prepare
	content := make(map[string]interface{})
	prepareMessage := &Message{
		Id:       -1,
		Action:   MessageActionPrepare,
		Content:  content,
		Callback: "", // 消费成功后 调用的回调函数
	}
	_, err := r.SendMessage(prepareMessage, exchange, routeKey)

	if err != nil {
		Logger.Error(prepareMessage.Action, err.Error())
		// panic(err)
	}
	return err
}

/**
事务开启 - 尝试与队列通信
@param exchange string 交换器名
@param routeKey string 路由key
@return error
*/
func TryMessageTransactionWithExchangeAndRoute(exchange string, routeKey string) error {
	return RabbitMq.TryMessageTransactionWithExchangeAndRoute(exchange, routeKey)
}
