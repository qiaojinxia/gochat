package repositories

import (
	"github.com/gomodule/redigo/redis"
	"goChat/conf"
	"sync"
	"time"
)

var redisPool *RedisServer
var initRedis = &sync.Once{}

func GetRedisPool() *redis.Pool {
	initRedis.Do(func() {
		server := RedisServer{}
		server.Init()
	})
	return redisPool.pool
}


type RedisServer struct {
	pool *redis.Pool
}

func NewRedisServer(pool *redis.Pool) *RedisServer {
	return &RedisServer{pool: pool}
}

func (rs *RedisServer) Close() error {
	return rs.pool.Close()
}

func (rs *RedisServer) Init() error {
	setDb := redis.DialDatabase(conf.RedisIndex)   //输入数据库序号
	setPasswd := redis.DialPassword(conf.RedisPwd) //设置密码

	pool := &redis.Pool{
		// maximum number of connections allocated by the pool at a given time.
		// when zero, there is no limit on the number of connections in the pool.
		//最大活跃连接数，0代表无限
		MaxActive: 888,
		//最大闲置连接数
		// maximum number of idle connections in the pool.
		MaxIdle: 20,
		//闲置连接的超时时间
		// close connections after remaining idle for this duration. if the value
		// is zero, then idle connections are not closed. applications should set
		// the timeout to a value less than the server's timeout.
		IdleTimeout: time.Second * 10,
		//定义拨号获得连接的函数
		// dial is an application supplied function for creating and configuring a
		// connection.
		//
		// the connection returned from dial must not be in a special state
		// (subscribed to pubsub channel, transaction started, ...).
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.RedisAddr, setDb, setPasswd)
		},
	}
	redisPool = NewRedisServer(pool)
	RegisterServer(redisPool)
	return nil
}
