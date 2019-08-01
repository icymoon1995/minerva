package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
	"strings"
)

// 对外提供 common.DB
var DB *xorm.Engine

func dbInit () {

	/**
	 *	初始化 数据库连接信息
	 */
	var err error

	// 获取环境
	var env = viper.GetString("common.env")

	// 默认为develop
	if strings.EqualFold(env, "") {
		env = "develop"
	}

	// 配置前缀 未使用全局的前缀 viper.SetEnvPrefix()
	var dbPrefix string = env + ".db."

	// 读取数据库连接的配置
	var username string = viper.GetString(dbPrefix + "username")
	var password string = viper.GetString(dbPrefix + "password")
	var host string = viper.GetString(dbPrefix + "host")
	var port string = viper.GetString(dbPrefix + "port")
	var database string = viper.GetString(dbPrefix + "database")
	var charset string = viper.GetString(dbPrefix + "charset")
	// 生成mysql的连接语句
	var databaseConfig string = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s", username, password, host, port, database, charset)

	// 连接xorm
	DB, err = xorm.NewEngine("mysql", databaseConfig)

	// 测试数据库连接是否 OK
	if err = DB.Ping(); err != nil {
		fmt.Println("ping db error:", err)
	}

	// 是否将 生成的sql语句打印在控制台
	var showSql bool = viper.GetBool(dbPrefix + "showsSql")
	DB.ShowSQL(showSql)

}