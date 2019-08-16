package main

import (
	"github.com/spf13/viper"
	"minerva/src/common"
	"minerva/src/routes"
	"strings"
)

func main() {

	// 部分初始化工作
	common.Init()

	// 注册路由
	routes.RegisterRoutes()

	// 监听 开始服务
	startListen()
}

// echo的开启监听
func startListen() {
	var address string = getAddress()
	routes.Route.Logger.Fatal(routes.Route.Start(address))
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
