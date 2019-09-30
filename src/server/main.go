package main

import (
	"github.com/micro/go-micro/web"
	"github.com/spf13/viper"
	"log"
	"minerva/src/common"
	"minerva/src/routes"
	"strings"
)

func main() {

	// 部分初始化工作
	common.Init()

	// 注册路由
	// routes.RegisterRoutes()

	// 队列接受消息
	// go logic.ReceiveMessage()

	// 监听 开始服务
	// startListen()

	//session := common.DB.NewSession()
	//defer session.Close()
	//err := session.Begin()
	//
	//if err != nil {
	//
	//}
	//user := &model.User{}
	//
	//fmt.Println("main before get user")
	//result,err := session.Where("id = ?" , 4).ForUpdate().Get(user)
	//fmt.Println("main after get user")
	//fmt.Println(result)
	//user.Name = "111112222"
	//
	//re, err := session.Where("id = ?", user.Id).Update(user)
	//
	//time.Sleep(5 * time.Second)
	//
	//fmt.Println(re)
	//
	//_ = session.Commit()

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
}

// echo的开启监听
func startListen() {

	// 获取监听地址
	var address string = getAddress()

	// 服务名
	var fullServiceName string = getServiceName()

	// go-micro 服务注册
	service := web.NewService(
		web.Name(fullServiceName),
		web.Address(address),
	)

	// 路由处理
	service.Handle("/", routes.Route)

	// go-micro 服务初始化
	err := service.Init()

	if err != nil {
		log.Fatal(err)
	}

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// echo的启动方式
	// routes.Route.Logger.Fatal(routes.Route.Start(address))
}

// 获取地址:端口的 字符串拼接
func getAddress() string {
	var host string = viper.GetString("common.host")
	var port string = viper.GetString("common.port")
	if strings.EqualFold(host, "") {
		host = "127.0.0.1"
	}
	if strings.EqualFold(port, "") {
		port = "8088"
	}
	return host + ":" + port
}

// 获取 服务名
func getServiceName() string {
	// 服务名
	var serviceName string = viper.GetString("common.serviceName")
	// 服务namespace
	var namespace string = viper.GetString("common.namespace")

	return namespace + "." + serviceName
}
