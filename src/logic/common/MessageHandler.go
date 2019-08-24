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

			// 处理逻辑
			handle(d.Body)

			log.Println("end handle message")

			d.Reject(false)
			if error != nil {
				_ = d.Reject(false)
			}
			_ = d.Ack(false)
			// 批量
			//_ = d.Ack(false)
			//
			// d.Ack(true)
			//err := d.Reject(true)
			//if err != nil {
			//	log.Print(err)
			//}
		}
	}()
	<-forever
}

func handle(body []byte) {

	// 默认message是key=>value格式
	message := make(map[string]string)
	error := json.Unmarshal(body, &message)

	if error != nil {
		log.Println("json.unmarshal error : ", error)
	}

	// 具体处理
	for key, value := range message {
		log.Printf(key + " : " + value)
	}
}
