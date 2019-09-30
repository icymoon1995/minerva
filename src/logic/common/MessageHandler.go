package logic

import (
	"encoding/json"
	"log"
	"minerva/src/common"
)

func ReceiveMessage() {
	msgs, error := common.Channel.Consume(
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
		for d := range msgs {

			if error != nil {
				log.Println("some error", error)
			}

			log.Println("start handle message")

			// d.MessageId
			// 处理逻辑
			error = handle(d.Body, d.DeliveryTag)

			log.Println("end handle message")

			//	d.Reject(false)
			if error != nil {
				_ = d.Reject(false)
			}
			_ = d.Ack(false)

		}
	}()
	<-forever
}

/**
	接收消息体 处理消息
 	@param body []byte 消息体内容
	@param deliveryTag uint64 消息编号tag
*/
func handle(body []byte, deliveryTag uint64) error {
	// 用interface处理消息
	// var jsonInterface interface{}
	// json.Unmarshal(d.Body, &jsonInterface)
	// message2 := jsonInterface.(map[string]interface{})
	// message3["id"], ":", message3["action"], ":", message3["callback"], ":", message3["content"]

	var message common.Message
	error := json.Unmarshal(body, &message)

	if error != nil {
		log.Println("json.unmarshal error : ", error)
	}

	// 具体的处理逻辑
	// 保证幂等性 用redis或者其他做唯一处理  redis.setnx("message:id")
	log.Println(message.Id, ":", message.Action, ":", message.Content, ":", message.Callback)

	return error
}
