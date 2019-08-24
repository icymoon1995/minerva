package common

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"net"
	"strings"
	"time"
)

// 全局的rabbitMQ
var rabbit *amqp.Connection
var Channel *amqp.Channel
var Queue amqp.Queue
var rabbitPrefix string // redis配置前缀

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

	// 创建通道
	createChannel()

	// 创建队列
	createQueue()

	// 绑定队列
	queueBind()

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
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	log.Println("rabbit mq connection success")
}

/**
创建通道
*/
func createChannel() {
	var err error
	Channel, err = rabbit.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open channel", err)
	}
	log.Println("rabbit mq channel success")
}

/**
创建队列
*/
func createQueue() {
	var name string = viper.GetString(rabbitPrefix + "queue")
	var err error

	args := amqp.Table{}
	// 绑定该队列的 死信route和key
	args["x-dead-letter-exchange"] = "dead_message_exchange"
	args["x-dead-letter-routing-key"] = "dead_key"

	Queue, err = Channel.QueueDeclare(
		name,  // name
		true,  // durable -- 是否持久化
		false, // delete when unused
		false, // exclusive -- 是否独占
		false, // no-wait -- 阻塞消息
		args,  // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to init queue", err)
	}
}

/**
队列绑定
*/
func queueBind() {
	// 绑定的route key
	var routeKeys string = viper.GetString(rabbitPrefix + "routekey")
	// 绑定的 交换机
	var exchangeName string = viper.GetString(rabbitPrefix + "exchange")
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
			log.Fatalf("%s: %s", "Failed to bind queue", err)
		}
	}

}

// 发送消息需要创建exchange
//func createExchange(name string, kind string) {
//	error := channel.ExchangeDeclare(
//		name, // name
//		kind,      // type
//		true,          // durable
//		false,         // auto-deleted
//		false,         // internal
//		false,         // no-wait
//		nil,           // arguments
//	)
//
//	if error != nil {
//		log.Fatalf("%s: %s", "Failed to create exchange", error)
//	}
//}

//func SendMessage(message map[string]interface{}, exchangeName string, routeName string) {
//	jsonMsg, _ := json.Marshal(message)
//
//	error := channel.Publish(
//		exchangeName, // exchange
//		routeName,    // routing key
//		false,        // mandatory
//		false,
//		amqp.Publishing{
//			DeliveryMode: amqp.Persistent,
//			ContentType:  "text/plain",
//			Body:         jsonMsg,
//			//[]byte(body),
//		})
//
//	if error != nil {
//		log.Fatalf("%s: %s", "Failed to send message", error)
//	}
//}
