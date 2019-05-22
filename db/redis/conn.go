package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
)

func newRedisPool() *redis.Pool {
	return &redis.Pool{
		// 最多有多少条连接
		MaxIdle: 100,
		// 最大有多少活动链接
		MaxActive: 80,
		// 链接超时时间
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 打开链接
			conn, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Printf("Failed to Dial redis,err:%s\n", err.Error())
				return nil, err
			}
			// 访问认证
			//if _,err = conn.Do("AUTH","");err != nil {
			//	conn.Close()
			//	fmt.Printf("Failed to auth redis,err:%s\n",err.Error())
			//	return nil,err
			//}
			return conn, nil
		},
		// 链接后测试
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func init() {
	pool = newRedisPool()
}

// 对外暴露连接池
func RedisPool() *redis.Pool {
	return pool
}
