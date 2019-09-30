package common

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
)

/**
	redigo
     对应redis的常用操作
*/
/**
 * 运行命令:
 * go test -v src/common/redigo_test.go
 * go test -v src/common/redigo_test.go -test.run {methodName}
 */
var redisClient redis.Conn

// 初始化连接 单连接
func initClient() redis.Conn {
	connectionType := "tcp"
	address := "127.0.0.1:6379"
	password := "root"
	database := 0
	redisClient, _ := redis.Dial(connectionType, address, redis.DialPassword(password), redis.DialDatabase(database))

	return redisClient
}

// 常用操作 - String GET
func TestGetString(t *testing.T) {
	redisClient := initClient()

	var notExistKey string = "key" // 无此key
	var existKey string = "hello"  // hello -> world
	// GET 方法默认返回的类型是 []uint8/[]byte
	// 由于有nil的关系 直接使用 string(reply.([]uint8)) 会有问题
	replyNotExist, error := redisClient.Do("GET", notExistKey)

	replyExist, error := redisClient.Do("GET", existKey)

	// exist, error := redis.Bool(redisClient.Do("EXISTS", existKey))

	// 也可以这样直接用
	// reply, _ := redis.String(redisClient.Do("GET", existKey))
	if error != nil {
		t.Error("redigo_test#GetStringTest, redisClient.Do GET error :", error)
	}

	t.Log(redis.String(replyNotExist, error)) // nil
	t.Log(redis.String(replyExist, error))    // world
}

// 常用操作 - String SET
func TestSetString(t *testing.T) {
	redisClient := initClient()

	var key string = "hey2"
	var value string = "world2"
	var expireType string = "EX" // EX/PX
	var timeout int = 30
	result, error := redisClient.Do("SET", key, value, expireType, timeout)

	// result, error := redisClient.Do("SETNX", key, value)
	// result, error := redisClient.Do("EXPIRE", key, timeout)

	if error != nil {
		t.Error("redigo_test#GetStringTest, redisClient.Do GET error :", error)
	}

	//if result.(int64) == 0 {
	//	fmt.Println(123)
	//} else {
	//	fmt.Println(5555)
	//}
	fmt.Println(result.(int64))
	t.Log(result) // OK
}

// 常用操作 - String DEL
func TestDelString(t *testing.T) {
	redisClient := initClient()

	var key string = "hey2"
	result, error := redisClient.Do("DEL", key)

	// result, error := redisClient.Do("SETNX", key, value)
	// result, error := redisClient.Do("EXPIRE", key, timeout)

	if error != nil {
		t.Error("redigo_test#TestDelString, redisClient.Do GET error :", error)
	}

	t.Log(result) // 1
}

// ttl 或者pttl 。
func TestTTL(t *testing.T) {
	redisClient := initClient()
	var key string = "hey2"
	// redis.Int(redisClient.Do("ttl", key))
	timeout, error := redisClient.Do("TTL", key)

	if error != nil {
		t.Error("redigo_test#TestTTL, redisClient.Do GET error :", error)
	}

	t.Log(timeout) // 时间。。int类型
}

// incr/decr/incrby/decrby 。
func TestChange(t *testing.T) {
	redisClient := initClient()
	var key string = "heyNum"
	var addValue int = 5
	result, error := redisClient.Do("INCRBY", key, addValue)

	if error != nil {
		t.Error("redigo_test#TestChange, redisClient.Do GET error :", error)
	}

	t.Log(result) // 返回更改后的结果
}

// keys
func TestKeys(t *testing.T) {
	redisClient := initClient()
	var selectKey string = "*"
	result, error := redis.Strings(redisClient.Do("KEYS", selectKey))

	if error != nil {
		t.Error("redigo_test#Keys, redisClient.Do GET error :", error)
	}

	t.Log(result) // 返回更改后的结果
}

// 常用操作 - hash set
func TestSetHash(t *testing.T) {
	redisClient := initClient()

	var key string = "user"
	var attributeKey1 string = "name"
	var attributeValue1 string = "user_name"

	var attributeKey2 string = "email"
	var attributeValue2 string = "123456@email.com"
	result, error := redisClient.Do("HMSET", key, attributeKey1, attributeValue1, attributeKey2, attributeValue2)

	if error != nil {
		t.Error("redigo_test#TestSetHash, redisClient.Do GET error :", error)
	}

	t.Log(result) // OK
}

// 常用操作 - Hash GET
func TestGetHash(t *testing.T) {
	redisClient := initClient()

	var key string = "user"
	var attributeKey1 string = "name"
	var attributeKey2 string = "email"
	// HGETALL
	result, error := redis.Strings(redisClient.Do("HMGET", key, attributeKey1, attributeKey2))

	if error != nil {
		t.Error("redigo_test#TestSetHash, redisClient.Do GET error :", error)
	}

	t.Log(result) // [value1, value2..]
}

// 常用操作 - list Set
func TestSetList(t *testing.T) {
	redisClient := initClient()

	var key string = "list"
	var attributeKey1 string = "attribute"
	var attributeKey2 string = "email"
	// LPUSH/RPUSH
	result, error := redisClient.Do("LPUSH", key, attributeKey1, attributeKey2)

	if error != nil {
		t.Error("redigo_test#TestSetList, redisClient.Do GET error :", error)
	}

	t.Log(result) // 2
}

// 常用操作 - list Get
func TestGetList(t *testing.T) {
	redisClient := initClient()

	var key string = "list"
	// LPOP/RPOP
	result, error := redis.String(redisClient.Do("RPOP", key))

	if error != nil {
		t.Error("redigo_test#TestGetList, redisClient.Do GET error :", error)
	}

	t.Log(result) // 2
}

// 常用操作 - list range
func TestRangeList(t *testing.T) {
	redisClient := initClient()

	var key string = "list"
	var start int = 0
	var end int = 5
	result, error := redis.Strings(redisClient.Do("LRANGE", key, start, end))

	if error != nil {
		t.Error("redigo_test#TestRangeList, redisClient.Do GET error :", error)
	}

	t.Log(result) // [...]
}

// 常用操作 - set set
func TestSetSet(t *testing.T) {
	redisClient := initClient()

	var key string = "set"
	var item1 string = "hello"
	var item2 string = "world"
	result, error := redisClient.Do("SADD", key, item1, item2)

	if error != nil {
		t.Error("redigo_test#TestRangeList, redisClient.Do GET error :", error)
	}

	t.Log(result) // [...]
}

//  常用操作 - set get all values 查看集合所有values
func TestValuesSet(t *testing.T) {
	redisClient := initClient()

	var key string = "set"
	result, error := redis.Strings(redisClient.Do("SMEMBERS", key))

	if error != nil {
		t.Error("redigo_test#TestValuesSet, redisClient.Do GET error :", error)
	}

	t.Log(result) // [...]
}

// 常用操作 - set sscan
func TestIteratorSet(t *testing.T) {
	redisClient := initClient()

	var (
		cursor int64
		items  []string
	)

	results := make([]string, 0)

	var key string = "set"
	for {
		values, error := redis.Values(redisClient.Do("SSCAN", key, cursor))

		if error != nil {
			t.Error("redigo_test#TestIteratorSet, redisClient.Do GET error :", error)
		}

		values, error = redis.Scan(values, &cursor, &items)
		if error != nil {
			t.Error("redigo_test#TestIteratorSet, redisClient.Scan GET error :", error)
		}
		results = append(results, items...)

		if cursor == 0 {
			break
		}
	}

	t.Log(results) // [...]
}

// 常用操作 - set key exist
func TestKeyExistSet(t *testing.T) {
	redisClient := initClient()

	var key string = "set"
	var validateKey string = "hello"
	result, error := redisClient.Do("SISMEMBER", key, validateKey)

	if error != nil {
		t.Error("redigo_test#TestIteratorSet, redisClient.Do GET error :", error)
	}

	t.Log(result) // 1/0
}

// 常用操作 - sets operate
func TestOperateSet(t *testing.T) {
	redisClient := initClient()

	var set string = "set"
	var set2 string = "set2"
	// SINTER 交集
	// SUNION 并集
	// SDIFF 差集
	result, error := redis.Strings(redisClient.Do("SUNION", set, set2))

	if error != nil {
		t.Error("redigo_test#TestOperateSet, redisClient.Do GET error :", error)
	}

	t.Log(result) // 1/0
}

// 常用操作 - sets len
func TestLenSet(t *testing.T) {
	redisClient := initClient()

	var set string = "set"
	// SINTER 交集
	// SUNION 并集
	result, error := redis.Int(redisClient.Do("SCARD", set))

	if error != nil {
		t.Error("redigo_test#TestLenSet, redisClient.Do GET error :", error)
	}

	t.Log(result) // int
}

// 常用操作 - sortset set
func TestSetSortSet(t *testing.T) {
	redisClient := initClient()

	var key string = "sortset"
	var score int = 1
	var item string = "world"
	result, error := redisClient.Do("ZADD", key, score, item)

	if error != nil {
		t.Error("redigo_test#TestSortSetSet, redisClient.Do GET error :", error)
	}

	t.Log(result) // [...]
}

// 常用操作 - sortset incr
func TestIncrSortSet(t *testing.T) {
	redisClient := initClient()

	var key string = "sortset"
	var addScore int = 3
	var item string = "world"
	result, error := redis.Int(redisClient.Do("ZINCRBY", key, addScore, item))

	if error != nil {
		t.Error("redigo_test#TestSortSetSet, redisClient.Do GET error :", error)
	}

	t.Log(result) // int
}

// 常用操作 - sortset zrank/zrevrank
func TestRankSortSet(t *testing.T) {
	redisClient := initClient()

	var key string = "sortset"
	var addScore int = 3
	var item string = "world"
	// zrank 从小到大
	// zrevrank 从大到小
	result, error := redis.Int(redisClient.Do("ZRANK", key, addScore, item))

	if error != nil {
		t.Error("redigo_test#TestSortSetRank, redisClient.Do GET error :", error)
	}

	t.Log(result) // int
}

// 常用操作 - sortset zrange/zrevrange
func TestRangeSortSet(t *testing.T) {
	redisClient := initClient()

	var key string = "sortset"
	var start int = 0
	var end int = 1
	// zrange 从小到大
	// zrevrange 从大到小
	result, error := redis.Strings(redisClient.Do("ZRANGE", key, start, end, "withscores"))

	if error != nil {
		t.Error("redigo_test#TestSortSetRange, redisClient.Do GET error :", error)
	}

	t.Log(result) // int
}

// 常用操作 - sortset zrem
func TestREMSortSet(t *testing.T) {
	redisClient := initClient()

	var key string = "sortset"
	var item string = "hello"
	result, error := redis.Int(redisClient.Do("zrem", key, item))

	if error != nil {
		t.Error("redigo_test#TestSortSetREM, redisClient.Do GET error :", error)
	}

	t.Log(result) // int
}

// 常用操作 - sortset zscan
func TestIteratorSortSet(t *testing.T) {
	redisClient := initClient()

	var (
		cursor int64
		items  []string
	)

	results := make([]string, 0)

	var key string = "sortset"
	for {
		values, error := redis.Values(redisClient.Do("ZSCAN", key, cursor))

		if error != nil {
			t.Error("redigo_test#TestIteratorSortSet, redisClient.Do GET error :", error)
		}

		values, error = redis.Scan(values, &cursor, &items)
		t.Log("items", items)

		if error != nil {
			t.Error("redigo_test#TestIteratorSortSet, redisClient.Scan GET error :", error)
		}
		results = append(results, items...)

		if cursor == 0 {
			break
		}
	}

	t.Log(results) // [...]
}
