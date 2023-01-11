package db

import (
	"errors"
	goredis "github.com/go-redis/redis/v8"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/irvankadhafi/erajaya-product-service/internal/config"
	"time"
)

// NewRedisConnPool uses redigo library to establish the redis connection pool
func NewRedisConnPool(url string) (*redigo.Pool, error) {
	_, err := goredis.ParseURL(url)
	if err != nil {
		return nil, errors.New("invalid redis URL: " + url)
	}

	return &redigo.Pool{
		MaxIdle:     config.RedisMaxIdleConn(),
		MaxActive:   config.RedisMaxActiveConn(),
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		MaxConnLifetime: 1 * time.Minute,
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Wait: true,
	}, nil
}
