package common

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var Enviorment string
var JWTKey string

/**
  处理配置文件信息
*/
func configInit() {

	// 处理ini 可以用 https://github.com/go-gcfg/gcfg

	// 配置文件类型
	viper.SetConfigType("yaml")
	// 配置文件路径
	viper.AddConfigPath("config/")
	// 配置文件名
	viper.SetConfigName("env")
	// 读取文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

func envInit() {
	// 获取环境
	var env = viper.GetString("common.env")

	// 默认为develop
	if strings.EqualFold(env, "") {
		env = "develop"
	}

	Enviorment = env

	JWTKey = viper.GetString(Enviorment + ".jwt.key")
}
