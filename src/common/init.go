package common

func Init() {
	// 日志初始化
	loggerInit()

	// 配置文件初始化
	configInit()

	// env变量初始化
	envInit()

	// 数据库连接初始化
	dbInit()

	// redis连接池初始化
	redisPoolInit()

	// rabbitMq初始化
	rabbitInit()

}
