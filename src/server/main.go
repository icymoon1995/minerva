package main

import (
	"github.com/micro/go-micro/web"
	"github.com/spf13/viper"
	"log"
	"minerva/src/common"
	logic "minerva/src/logic/common"
	"minerva/src/routes"
	"strings"
)

func main() {

	// 部分初始化工作
	common.Init()

	// 注册路由
	routes.RegisterRoutes()

	// 队列接受消息
	go logic.ReceiveMessage()

	// 监听 开始服务
	startListen()
}

type Listener1 struct {
	Status int
}

func (listener *Listener1) Check() int {

	return listener.Status
}

func (listener *Listener1) RollBack() error {
	log.Println("listener rollback")
	listener.Status = 2
	return nil
}

func (listener *Listener1) Execute() error {
	log.Println("execute local transaction")
	listener.Status = 1
	return nil
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
