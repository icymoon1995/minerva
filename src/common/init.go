package common

import (
	"math/rand"
	"time"
)

// 全局配置
var Global = &config{}

func Init() {
	// 日志初始化
	loggerInit()

	// 随机数种子 -- 进程级别 保证每个进程的seed都不一致
	rand.Seed(time.Now().UnixNano())

	// 配置文件初始化
	configInit()

	// 数据库连接初始化
	newDB()

	// redis连接池初始化
	newRedisPool()

	// rabbitMq初始化
	newRabbit()

}
