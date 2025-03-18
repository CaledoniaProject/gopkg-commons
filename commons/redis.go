package commons

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewRedisPool(dsn string) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(dsn)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return nil, err
	}

	return pool, nil
}
