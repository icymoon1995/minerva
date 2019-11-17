package common

import (
	"github.com/spf13/viper"
)

/**
使用事务的小demo:

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

main() {
	listen := Listener1{

	}
	trans := &common.NewTransaction(&listen)

	message := common.Message{
		Id:1,
		Action:"add",
		Content:map[string]interface{}{
			"message" : "message content test",
		},
		Callback:"callback_url",
	}

	err := trans.MakeMessageTransaction(message)

	fmt.Println(err)
}
*/

// 事务监听接口
type TransactionListener interface {
	Execute() error  // 执行事务方法
	Check() int      // 查询方法 主要用于 消息队列无法通信时作为回调查询
	RollBack() error // 回滚方法
}

//
type transaction struct {
	ExchangeName string              // 专门处理事务的交换器
	RouteKey     string              // 专门处理事务的路由key
	Listener     TransactionListener // 对应的listener
	NeedTry      bool                // 是否发送try消息 确认存活 默认关闭
}

/**
rabbitMq 事务类初始化
*/
func NewTransaction(listener TransactionListener) *transaction {

	trans := new(transaction)
	trans.ExchangeName = viper.GetString(rabbitPrefix + "transaction.exchangeName")
	trans.RouteKey = viper.GetString(rabbitPrefix + "transaction.routeKey")
	trans.Listener = listener
	trans.NeedTry = true

	var exchangeType string = viper.GetString(rabbitPrefix + "transaction.exchangeType")
	var queueName string = viper.GetString(rabbitPrefix + "transaction.queueName")

	// 创建exchange
	go exchangeInit(trans.ExchangeName, exchangeType)

	// 创建队列
	go queueInit(queueName, "", "")

	// 队列绑定
	go queueBind(queueName, trans.RouteKey, trans.ExchangeName)

	return trans
}

/**
发送一条ping消息
@return error
*/
func (trans *transaction) Try() error {
	// 走默认的
	return TryMessageTransactionWithExchangeAndRoute(trans.ExchangeName, trans.RouteKey)
}

/**
发送事务处理的消息
@param trulyMessage Message 消息
@return error
*/
func (trans *transaction) MakeMessageTransaction(trulyMessage *Message) error {
	var err error

	// 如果需要  发送try 消息
	if trans.NeedTry {
		err = trans.Try()
		if err != nil {
			return err
		}
	}

	// 执行事务
	err = trans.Listener.Execute()
	if err != nil {
		return err
	}

	// 发送消息
	confirmation, err := SendMessage(trulyMessage, trans.ExchangeName, trans.RouteKey)

	// 出现异常、 ||  ack = false （拒绝消息)
	if err != nil || (confirmation != nil && confirmation.Ack == false) {
		Logger.Println("transactionWithRabbit #MakeMessageTransaction error:", err.Error())
		// trulyMessage.Action, err.Error())
		var tryTime int = 0
		// 回滚操作
		for {
			// 出现异常则进行回滚 直到回滚成功
			err = trans.Listener.RollBack()
			if err == nil {
				break
			}
			if tryTime >= MaxTries {
				Logger.Println("transactionWithRabbit #MakeMessageTransaction rollback err : ", err.Error())
				// 其他报警机制
			}
			tryTime++

		}
	}
	return err
}
