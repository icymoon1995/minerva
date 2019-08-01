package common

import (
	"fmt"
	"github.com/spf13/viper"
)

/**
  处理配置文件信息
*/
func configInit () {

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
