package global

import (
	"github.com/spf13/viper"
	"runtime"
	"time"
)

var Global = &config{}

type config struct {
	Environment string
	JWTKey      string
	GoVersion   string
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
	c.JWTKey = viper.GetString(c.Environment + ".jwt.key")
	return nil
}
