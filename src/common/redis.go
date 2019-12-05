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

var redisPool = &redisPoolConfig{}

type redisPoolConfig struct {
	prefix      string
	maxActive   int
	maxIdle     int
	idleTimeout time.Duration

	redisConfig struct {
		// 读取redis连接的配置
		// 连接方式 默认为tcp
		connectionType string
		// 地址
		host string
		// Auth密码
		password string
		// 端口
		port string
		// 使用数据库
		database int
		// 连接超时时间
		connTimeout time.Duration
		// 读超时时间
		readTimeout time.Duration
		// 写超时时间
		writeTimeout time.Duration
	}
}

func newRedisPool() {
	// 初始化前缀
	redisPool.redisPoolInit()

	RedisPool = &redis.Pool{

		MaxIdle:     redisPool.maxIdle,     //  最大空闲线程数
		MaxActive:   redisPool.maxActive,   // 最大活跃线程数
		IdleTimeout: redisPool.idleTimeout, // 空闲等待连接时间

		Dial: func() (redis.Conn, error) {
			client, err := redisPool.redisConn()
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
redis连接池初始化
*/
func (redisPool *redisPoolConfig) redisPoolInit() {
	// 初始化前缀
	redisPool.prefix = Global.Environment + ".redis."
	redisPool.maxActive = viper.GetInt(redisPool.prefix + "maxActive")
	redisPool.maxIdle = viper.GetInt(redisPool.prefix + "maxIdle")
	redisPool.idleTimeout = 3 * time.Second

	// 设置default的也行。。
	var connectionType string = viper.GetString(redisPool.prefix + "connection")
	if strings.EqualFold(connectionType, "") {
		connectionType = "tcp"
	}
	var host string = viper.GetString(redisPool.prefix + "host") // 地址
	if strings.EqualFold(host, "") {
		host = "127.0.0.1"
	}

	redisPool.prefix = Global.Environment + ".redis."
	redisPool.redisConfig.host = host
	redisPool.redisConfig.connectionType = connectionType                           // 连接方式 默认为tcp
	redisPool.redisConfig.password = viper.GetString(redisPool.prefix + "password") // Auth密码
	redisPool.redisConfig.port = viper.GetString(redisPool.prefix + "port")         // 端口
	redisPool.redisConfig.database = viper.GetInt(redisPool.prefix + "database")    // 默认为0
	redisPool.redisConfig.connTimeout = 3 * time.Second
	redisPool.redisConfig.readTimeout = 3 * time.Second
	redisPool.redisConfig.writeTimeout = 3 * time.Second

}

/**
配置连接redis
*/
func (redisPool *redisPoolConfig) redisConn() (redis.Conn, error) {

	var address string = redisPool.redisConfig.host + ":" + redisPool.redisConfig.port

	redisClient, error := redis.Dial(
		redisPool.redisConfig.connectionType,
		address,
		redis.DialPassword(redisPool.redisConfig.password),
		redis.DialDatabase(redisPool.redisConfig.database),
		redis.DialConnectTimeout(3*time.Second),
		redis.DialReadTimeout(3*time.Second),
		redis.DialWriteTimeout(3*time.Second),
	)
	if error != nil {
		Logger.Fatal("redis.go #redisConn error", error)
	}

	return redisClient, error
}
