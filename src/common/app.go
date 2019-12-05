package common

import (
	"fmt"
	"github.com/spf13/viper"
	"runtime"
	"strconv"
	"time"
)

type config struct {
	Environment string
	JWTConfig   struct {
		Key    string
		Expire time.Duration
	}
	GoVersion string
	// 启动时间
	LaunchTime time.Time
}

func (c *config) initField() error {
	// 给默认值
	viper.SetDefault("common.env", "develop")
	/**
	// 默认为develop
	if strings.EqualFold(env, "") {
		env = "develop"
	}
	*/
	c.Environment = viper.GetString("common.env")
	c.GoVersion = runtime.Version()
	c.LaunchTime = time.Now()
	c.JWTConfig.Key = viper.GetString(c.Environment + ".jwt.key")

	// 过期时间 -- 必须转化成time.Duration格式 不然会抛异常
	var jwtExpire string = viper.GetString(c.Environment + ".jwt.expire")
	jwtExpireInt, _ := strconv.ParseInt(jwtExpire, 10, 64)
	// 目前单位是
	c.JWTConfig.Expire = time.Duration(jwtExpireInt) * time.Second

	return nil
}

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
		panic(fmt.Errorf("global init configInit viper readInConfig file: %s \n", err))
	}

	// 初始化 配置
	err = Global.initField()
	if err != nil {
		panic(fmt.Errorf("global init configInit app fillFields: %s \n", err))
	}

}
