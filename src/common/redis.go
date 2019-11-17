package common

import (
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// 全局的redis连接池  初始化后使用common.RedisPool.Get()即可
// 使用get请 务必带上 defer close()， 不然资源会一直占用, 导致无法释放
var RedisPool *redis.Pool

var prefix string // redis配置前缀

// 初始化前缀
func preInit() {
	prefix = Enviorment + ".redis."
}

/**
redis连接池初始化
*/
func redisPoolInit() {
	// 初始化前缀
	preInit()

	var maxActive int = viper.GetInt(prefix + "maxActive")
	var maxIdle int = viper.GetInt(prefix + "maxIdle")

	RedisPool = &redis.Pool{

		MaxIdle:     maxIdle,         //  最大空闲线程数
		MaxActive:   maxActive,       // 最大活跃线程数
		IdleTimeout: 3 * time.Second, // 空闲等待连接时间

		Dial: func() (redis.Conn, error) {
			client, err := redisConn()
			if err != nil {
				return nil, err
			}
			return client, err
		},

		Wait: true, // 队列处理策略 超过最大数(maxActive) 选择等待, 为false 超过maxActive会报 connection pool exhausted

		TestOnBorrow: func(c redis.Conn, t time.Time) error { // 连接心跳检测 -- 可能会影响性能
			_, err := c.Do("PING")
			if err != nil {
				Logger.Fatal("redis.go #redisPoolInit err:", err)
			}
			return err
		},
	}
}

/**
配置连接redis
*/
func redisConn() (redis.Conn, error) {
	// 配置前缀 未使用全局的前缀 viper.SetEnvPrefix()

	// 读取redis连接的配置
	// 连接方式 默认为tcp
	var connectionType string = viper.GetString(prefix + "connection")
	// 地址
	var host string = viper.GetString(prefix + "host")
	// Auth密码
	var password string = viper.GetString(prefix + "password")
	// 端口
	var port string = viper.GetString(prefix + "port")
	// 使用数据库
	var database int = viper.GetInt(prefix + "database") // 默认为0

	if strings.EqualFold(connectionType, "") {
		connectionType = "tcp"
	}
	if strings.EqualFold(host, "") {
		host = "127.0.0.1"
	}

	// databaseInt, _ = strconv.Atoi(database)
	// var connTimeout string = viper.GetString(dbPrefix + "conn_timeout")
	// var readTimeout string = viper.GetString(dbPrefix + "read_timeout")
	// var writeTimeout string = viper.GetString(dbPrefix + "write_timeout")

	var address string = host + ":" + port

	// fmt.Println(address)
	redisClient, error := redis.Dial(
		connectionType,
		address,
		redis.DialPassword(password),
		redis.DialDatabase(database),
		redis.DialConnectTimeout(3*time.Second),
		redis.DialReadTimeout(3*time.Second),
		redis.DialWriteTimeout(3*time.Second),
	)
	if error != nil {
		Logger.Fatal("redis.go #redisConn error", error)
	}

	return redisClient, error
}
