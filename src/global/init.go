package global

import (
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"sync"
	"time"
)

var once = new(sync.Once)

/**
  处理配置文件信息
*/
func configInit() {

	once.Do(func() {

		// 随机数种子 -- 进程级别 保证每个进程的seed都不一致
		rand.Seed(time.Now().UnixNano())

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
			panic(fmt.Errorf("global init configInit viper readInConfig file: %s \n", err))
		}

		// 初始化 配置
		err = Global.initField()
		if err != nil {
			panic(fmt.Errorf("global init configInit app fillFields: %s \n", err))
		}
	})

}
