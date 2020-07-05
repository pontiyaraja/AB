package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// redisConnURLMap connection URL map
type redisConn struct {
	connection *redis.Client
	mutex      *sync.RWMutex
}

func (c *redisConn) setConnection(redisConnection *redis.Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.connection = redisConnection
}

func (c *redisConn) getConnection() (*redis.Client, error) {
	c.mutex.RLock()
	//defer c.mutex.RUnlock()
	// check if connection is cached
	redisConnection := c.connection
	// if not cached 1.read redis server url from vault 2. connect to redis 3. cache the connection
	c.mutex.RUnlock()
	if redisConnection == nil { //1.  not cached
		redisClient := connectRedis("localhost:6379", "", 0)
		c.setConnection(redisClient)
		return redisClient, nil
	}
	return redisConnection, nil
}

var (
	connectionList *redisConn
	rOnce          sync.Once
)

func newConectionList() *redisConn {

	rOnce.Do(func() {
		connectionList = &redisConn{
			mutex: new(sync.RWMutex),
		}
	})
	return connectionList
}

func init() {
	connectionList = newConectionList()
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

//Ping is to ping a connection
func Ping() {
	ctx := context.Background()
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return
	}
	pong, err := redisClient.Ping(ctx).Result()
	fmt.Println(pong, err)

}

//Set sets the value for given key
func Set(key string, value string) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	err := redisClient.Set(ctx, key, value, 0).Err()
	return err
}

// SetEx is to set the key value with an expiry
func SetEx(key string, value string, expireSec int64) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return nil
	}
	ctx := context.Background()
	err := redisClient.SetNX(ctx, key, value, time.Duration(expireSec)*time.Second)
	return err.Err()
}

// Get gets the value for given key
func Get(key string) (string, error) {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return "", connectionError
	}

	ctx := context.Background()
	r := redisClient.Get(ctx, key)
	d, e := r.Result()
	return d, e
}

// RPush is to push an object to an array
func RPush(listName string, serializedObj string) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	e := redisClient.LPush(ctx, listName, serializedObj)
	return e.Err()
}

// Del is to delete a key value from the cache
func Del(key string) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	s := redisClient.Del(ctx, key)
	return s.Err()
}

// Expire is to set an expiry to a key
func Expire(key string, expiration int64) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	return redisClient.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
}

/*-------------------------- HASH MAP -----------------------*/

// HDel deletes a key or list of keys from hashmap
func HDel(key string, fields ...string) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	s := redisClient.HDel(ctx, key, fields...)
	return s.Err()
}

// HMset sets multiple key value pairs into a hashmap
func HMset(key string, fields map[string]string) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	s := redisClient.HMSet(ctx, key, fields)
	return s.Err()
}

// HSet sets a key value pair in a hashmap
func HSet(key, field string, value string) error {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return connectionError
	}
	ctx := context.Background()
	s := redisClient.HSet(ctx, key, field, value)
	return s.Err()
}

// HGet gets the given key and its value from a hashmap
func HGet(key, field string) (string, error) {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return "", connectionError
	}
	ctx := context.Background()
	s := redisClient.HGet(ctx, key, field)
	return s.Result()
}

//HgetAll - gets all the keys of a hash map
func HgetAll(key string) (map[string]string, error) {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return nil, connectionError
	}
	ctx := context.Background()
	return redisClient.HGetAll(ctx, key).Result()
}

// HExists - checks if key exits in a hashmap
func HExists(key, field string) bool {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return false
	}
	ctx := context.Background()
	return redisClient.HExists(ctx, key, field).Val()
}

// HIncrBy - Increments the value of a key by the given incr
func HIncrBy(key, field string, incr int64) int64 {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return 0 // Todo find a better method for returing value
	}
	ctx := context.Background()
	return redisClient.HIncrBy(ctx, key, field, incr).Val()
}

//TTL - time to live value for a key
func TTL(key string) (time.Duration, error) {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return 0, connectionError
	}
	ctx := context.Background()
	r := redisClient.TTL(ctx, key)
	d, e := r.Result()
	return d, e
}

//Exists checks if the key exists and returns the value
func Exists(key string) bool {
	redisClient, connectionError := connectionList.getConnection()
	if connectionError != nil {
		return false
	}
	return redisClient.Exists(key).Val()
}

func ScanKeys(pattern string, count int64) ([]string, error) {
	redisClient, err := connectionList.getConnection()
	if err != nil {
		return nil, err
	}

	var c uint64
	var keys []string
	for {
		var k []string
		k, c = redisClient.Scan(c, pattern, count).Val()
		keys = append(keys, k...)
		if c == 0 {
			break
		}
	}
	return keys, nil
}
