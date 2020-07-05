package core

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/pontiyaraja/AB/ablog"
)

func init() {

	connectionList = newConectionList()
}

var (
	connectionList *redisConn
	rOnce          sync.Once
)

type redisConn struct {
	connection *redis.Client
}

func (c *redisConn) setConnection(redisConnection *redis.Client) {
	c.connection = redisConnection
}
func (c *redisConn) getConnection() (*redis.Client, error) {
	redisConnection := c.connection
	// if not cached 1.read redis server url from vault 2. connect to redis 3. cache the connection
	if redisConnection == nil { //1.  not cached
		redisClient := connectRedis("localhost:6379", "", 0)
		c.setConnection(redisClient)
		return redisClient, nil
	}
	return redisConnection, nil
}

func newConectionList() *redisConn {
	rOnce.Do(func() {
		connectionList = &redisConn{}
	})
	return connectionList
}
func connectRedis(address, password string, db int) *redis.Client {
	return redis.NewClient(loadRedisOptions(address, password, db))
}
func loadRedisOptions(address, password string, db int) *redis.Options {
	return &redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	}
}

// func connect() *redis.Client {
// 	redisclient := redis.NewClient(&redis.Options{
// 		Addr:            "localhost:6379",
// 		Password:        "",
// 		DB:              0,
// 		MaxRetries:      5,
// 		MaxRetryBackoff: time.Duration(5 * time.Second),
// 		PoolTimeout:     time.Duration(3 * time.Second),
// 	})
// 	redisclient.Ping().Name()
// 	return redisclient
// }

func LPop(listName string) (string, error) {
	c, err := connectionList.getConnection()
	if err != nil {
		ablog.Error("", err, nil)
	}

	strCmd := c.LPop(listName)
	return strCmd.Result()
	//fmt.Println(res, err)
}

func RPush(listName, value string) (int64, error) {
	c, err := connectionList.getConnection()
	if err != nil {
		ablog.Error("", err, nil)
	}
	strCmd := c.RPush(listName, value)
	return strCmd.Result()
	//fmt.Println(res, err)
}

func HSet(listName, apikey string, timeval string) (bool, error) {
	c, err := connectionList.getConnection()
	if err != nil {
		ablog.Error("", err, nil)
	}
	boolCmd := c.HSet(listName, apikey, timeval)
	return boolCmd.Result()
	//fmt.Println(res, err)
}

func HGet(listName, apikey string) (string, error) {
	c, err := connectionList.getConnection()
	if err != nil {
		ablog.Error("", err, nil)
	}
	strCmd := c.HGet(listName, apikey)
	return strCmd.Result()
	//fmt.Println(res, err)
}

// func FlushAll() (string, error) {
// 	c, err := connectionList.getConnection()
// 	if err != nil {
// 		kiplog.KIPError("", err, nil)
// 	}
// 	statCmd := c.FlushAll()
// 	return statCmd.Result()
// 	//fmt.Println(res, err)
