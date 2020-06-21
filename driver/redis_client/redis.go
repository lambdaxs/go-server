package redis_client

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisDB struct {
	DSN         string        //"127.0.0.1:6379"
	Password    string        //""
	DB          int           // 0
	MaxIdle     int           //100
	MaxActive   int           //500
	IdleTimeout time.Duration //6min

	DialTimeout  time.Duration //500
	ReadTimeout  time.Duration //2s
	WriteTimeout time.Duration //3s
	KeepAlive    time.Duration //5min
}

func (r *RedisDB) ConnectRedisPool() (pool *redis.Pool) {
	if r.MaxIdle == 0 {
		r.MaxIdle = 100
	}
	if r.MaxActive == 0 {
		r.MaxActive = 500
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = time.Second * 480
	}
	if r.KeepAlive == 0 {
		r.KeepAlive = time.Minute * 5
	}
	if r.DialTimeout == 0 {
		r.DialTimeout = time.Second * 2
	}
	if r.ReadTimeout == 0 {
		r.ReadTimeout = time.Second * 2
	}
	if r.WriteTimeout == 0 {
		r.WriteTimeout = time.Second * 3
	}
	return &redis.Pool{
		MaxIdle:     r.MaxIdle,
		MaxActive:   r.MaxActive,
		IdleTimeout: r.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.DSN,
				redis.DialConnectTimeout(r.DialTimeout),
				redis.DialReadTimeout(r.ReadTimeout),
				redis.DialWriteTimeout(r.WriteTimeout),
				redis.DialPassword(r.Password),
				redis.DialDatabase(r.DB),
				redis.DialKeepAlive(r.KeepAlive),
			)
			if err != nil {
				return nil, err
			}
			if len(r.Password) > 0 {
				if _, err := c.Do("AUTH", r.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if r.DB > 0 && r.DB < 16 {
				if _, err := c.Do("SELECT", r.DB); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}
