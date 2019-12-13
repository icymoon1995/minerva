package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/**
 */

var db = &dbConfig{}

// 对外提供 common.DB
var DB *xorm.Engine

type dbConfig struct {
	prefix   string // 配置文件前缀
	username string
	password string
	host     string // int 可能更好一些
	port     string // int 可能更好一些
	database string
	charset  string
	showSql  bool // 是否将 生成的sql语句打印在控制台
}

func (db *dbConfig) dbInit() {
	// 配置前缀 未使用全局的前缀 viper.SetEnvPrefix()
	db.prefix = Global.Environment + ".db."
	// 读取数据库连接的配置
	db.username = viper.GetString(db.prefix + "username")
	db.password = viper.GetString(db.prefix + "password")
	db.host = viper.GetString(db.prefix + "host")
	db.port = viper.GetString(db.prefix + "port")
	db.database = viper.GetString(db.prefix + "database")
	db.charset = viper.GetString(db.prefix + "charset")
	db.showSql = viper.GetBool(db.prefix + "showsSql")
}

func (db *dbConfig) dataSource() string {
	return fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=%s",
		db.username,
		db.password,
		db.host,
		db.port,
		db.database,
		db.charset)

}

func newDB() {
	/**
	 *	初始化 数据库连接信息
	 */
	var err error

	db.dbInit()

	// 生成mysql的连接语句
	var dataSource string = db.dataSource()

	// 连接xorm
	DB, err = xorm.NewEngine("mysql", dataSource)

	// 测试数据库连接是否 OK
	if err = DB.Ping(); err != nil {

		Logger.WithFields(logrus.Fields{
			"file":   "db.go",
			"method": "newDB",
			"error":  "ping error",
		}).Fatalln(err)
		//log.al("ping db error:", err)
	}

	// 是否将 生成的sql语句打印在控制台
	DB.ShowSQL(db.showSql)

}
