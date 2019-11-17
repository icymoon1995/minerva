package logic

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"minerva/src/common"
	"net/http"
	"strconv"
	"strings"
)

func ReceiveMessage() {

	forever := make(chan bool)

	go func() {

		//if common.Channel.IsClosed() || ch.IsClosed() {
		//	log.Println("连接断开，重新连接")
		//	err = InitRabbitmq(Url)
		//	log.Println(err)
		//}

		messages, err := common.Channel.Consume(
			common.Queue.Name, // queue
			"",                // consumer
			false,             // auto-ack
			false,             // exclusive
			false,             // no-local
			false,             // no-wait
			nil,               // args
		)

		for d := range messages {
			fmt.Println(d)
			if err != nil {
				common.Logger.Errorln("MessageHandler #ReceiveMessage error: ", err)
			}

			common.Logger.Println("MessageHandler #ReceiveMessage start handle message")

			// 处理逻辑
			go handle(d)

			common.Logger.Println("end handle message")

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
		common.Logger.Errorln("MessageHandler #handle json.unmarshal error : ", err)
		_ = d.Reject(false)
	}

	session := common.DB.NewSession()
	fmt.Println(message.Id, message.Action, message.Content)

	redisClient := common.RedisPool.Get()
	var redisKey string = "message:" + strconv.Itoa(message.Id)

	var redisValue string = "processing"

	defer func() {
		if err := recover(); err != nil {
			common.Logger.Errorln("MessageHandler#handle error : ", err)

			//	_ = session.Rollback()

		}
		_, _ = redisClient.Do("del", redisKey)
		// 关闭redis连接
		redisClient.Close()
		// 关闭session
		session.Close()
	}()

	// 保证幂等性 用redis或者其他做唯一处理  redis.setnx("message:id")
	result, err := redisClient.Do("setnx", redisKey, redisValue)
	if err != nil { // 出现异常 或者 已经在执行 属于重复消息，则拒绝执行
		common.Logger.Errorln("MessageHandler#handle message repeat :", redisKey)
		//	panic(err)
		return
	}

	if result.(int64) == 0 {
		_ = d.Reject(false)
	}

	// 开启mysql事务
	err = session.Begin()
	if err != nil {
		common.Logger.Errorln("MessageHandler#handle begin session error :", err)
		panic(err)
	}

	// 具体的处理逻辑
	common.Logger.Println(message.Id, message.Action, message.Content)

	if err != nil {
		_ = d.Reject(false)
		panic(err)
	}

	// 确认消息回复   ack需要做重试机制？
	err = d.Ack(false)

	if err != nil {
		common.Logger.Errorln("MessageHandler #handle ack error:", err)
		panic(err)
	}

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

	//提交事务
	_ = session.Commit()

}
