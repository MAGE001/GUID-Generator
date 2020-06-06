package storage

import (
	"log"
	"time"

	"GUID-Generator/conf"
	"github.com/gomodule/redigo/redis"
)

// rStorage impl NodeIdStorager interface.
type rStorage struct {
	client *redis.Pool
}

func NewRStorage() *rStorage {
	return &rStorage{
		client: redisClient(),
	}
}

// NextNodeId get next node id from redis.
func (r *rStorage) NextNodeId() (int64, error) {
	conn := r.client.Get()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("conn close failed, err: %s\n", err)
		}
	}()
	return redis.Int64(conn.Do("INCR", conf.New().NodeIdKey))
}

func redisClient() *redis.Pool {
	return &redis.Pool{
		MaxIdle:         conf.New().Redis.PoolSize,
		IdleTimeout:     time.Duration(conf.New().Redis.IdleTimeout) * time.Second,
		MaxActive:       2 * conf.New().Redis.PoolSize,
		Wait:            true,
		MaxConnLifetime: time.Minute * 30,

		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				conf.New().Redis.Addr,
				redis.DialDatabase(conf.New().Redis.Db),
				redis.DialReadTimeout(time.Duration(conf.New().Redis.DialReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(conf.New().Redis.DialWriteTimeout)*time.Second),
				redis.DialKeepAlive(time.Duration(conf.New().Redis.DialKeepAlive)*time.Second),
			)
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Printf("redis ping failed, err: %s\n", err)
			}
			return err
		},
	}
}
