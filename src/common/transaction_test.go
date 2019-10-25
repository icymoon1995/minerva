package common

import (
	"fmt"
	"log"
	"testing"
)

/**
* 运行命令:
todo 补go mock
*/

// 初始化连接 单连接
func TestProducer(t *testing.T) {
	listen := Listener1{}
	trans := NewTransaction(&listen)

	message := Message{
		Id:     1,
		Action: "add",
		Content: map[string]interface{}{
			"message": "message content test",
		},
		Callback: "callback_url",
	}

	err := trans.MakeMessageTransaction(message)

	fmt.Println(err)
}

type Listener1 struct {
}

func (listener *Listener1) Check() int {
	return 1
}

func (listener *Listener1) RollBack() error {
	log.Println("listener rollback")
	return nil
}

func (listener *Listener1) Execute() error {
	log.Println("11111")
	return nil
}
