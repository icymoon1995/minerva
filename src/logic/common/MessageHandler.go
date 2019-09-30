package logic

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"minerva/src/common"
	"net/http"
	"strconv"
	"strings"
)

func ReceiveMessage() {
	messages, err := common.Channel.Consume(
		common.Queue.Name, // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	forever := make(chan bool)

	go func() {
		for d := range messages {

			if err != nil {
				log.Println("some error", err)
			}

			log.Println("start handle message")

			// d.MessageId
			// 处理逻辑
			handle(d)

			log.Println("end handle message")

		}
	}()
	<-forever
}

/**
	接收消息体 处理消息
 	@param body []byte 消息体内容
	@param deliveryTag uint64 消息编号tag
*/
func handle(d amqp.Delivery) {
	// body []byte, deliveryTag uint64) error {
	// 用interface处理消息
	// var jsonInterface interface{}
	// json.Unmarshal(d.Body, &jsonInterface)
	// message2 := jsonInterface.(map[string]interface{})
	// message3["id"], ":", message3["action"], ":", message3["callback"], ":", message3["content"]

	// d.Body, d.DeliveryTag
	var message common.Message
	err := json.Unmarshal(d.Body, &message)

	if err != nil {
		log.Println("json.unmarshal error : ", err)
		_ = d.Reject(false)
	}

	session := common.DB.NewSession()
	redisClient := common.RedisPool.Get()
	var redisKey string = "message:" + strconv.Itoa(message.Id)
	var redisValue string = "processing"

	defer func() {
		if err := recover(); err != nil {
			log.Println("MessageHandler#handle error : ", err)
			_, _ = redisClient.Do("del", redisKey)
			_ = session.Rollback()

		}
		// 关闭redis连接
		redisClient.Close()
		// 关闭session
		session.Close()
	}()

	// 保证幂等性 用redis或者其他做唯一处理  redis.setnx("message:id")
	result, err := redisClient.Do("setnx", redisKey, redisValue)
	if err != nil { // 出现异常 或者 已经在执行 属于重复消息，则拒绝执行
		panic(err)
	}

	if result.(int64) == 0 {
		_ = d.Reject(false)
	}

	// 开启mysql事务
	err = session.Begin()
	if err != nil {
		panic(err)
	}

	// 具体的处理逻辑
	log.Println(message.Id, message.Action, message.Content)

	if message.Callback != "" {
		// 调用callback
		// todo 改为内部rpc或者给外部http调用
		resp, err := http.Post(message.Callback, "application/x-www-form-urlencoded",
			strings.NewReader("messageId="+strconv.Itoa(message.Id)))
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
	}

	//if err != nil {
	//	_ = d.Reject(false)
	//	panic(err)
	//}

	// 确认消息回复
	_ = d.Ack(false)
	// 提交事务
	_ = session.Commit()
}
