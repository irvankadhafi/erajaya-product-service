package cache

import (
	redigo "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	defaultTTL    = 10 * time.Second
	defaultNilTTL = 5 * time.Minute
)

type Cache interface {
	Get(key string) (any, error)
	Store(item Item) error
	DeleteByKeys(keys []string) error
	StoreNil(cacheKey string) error
	SetDefaultTTL(time.Duration)
	SetConnectionPool(*redigo.Pool)
}

type cache struct {
	connPool   *redigo.Pool //disesuaikan dengan library yang digunakan
	defaultTTL time.Duration
	nilTTL     time.Duration
}

func NewCache() Cache {
	return &cache{
		defaultTTL: defaultTTL,
		nilTTL:     defaultNilTTL,
	}
}

// SetDefaultTTL sets the defaultTTL field of the cacheManager struct to the given duration.
// The defaultTTL field specifies the default time-to-live (TTL) for items in the cache.
func (c *cache) SetDefaultTTL(d time.Duration) {
	c.defaultTTL = d
}

func (c *cache) SetConnectionPool(conn *redigo.Pool) {
	c.connPool = conn
}

func (c *cache) Get(key string) (cachedItem any, err error) {
	cachedItem, err = c.get(key)
	logrus.Warn(cachedItem)
	if err != nil && err != ErrKeyNotExist && err != redigo.ErrNil || cachedItem != nil {
		return
	}

	return nil, nil
}

func (c *cache) Store(item Item) error {
	client := c.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	_, err := client.Do("SETEX", item.GetKey(), c.decideCacheTTL(item), item.GetValue())
	return err
}

// DeleteByKeys Delete by multiple keys
func (c *cache) DeleteByKeys(keys []string) error {
	if len(keys) <= 0 {
		return nil
	}

	client := c.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	var redisKeys []any
	for _, key := range keys {
		redisKeys = append(redisKeys, key)
	}

	_, err := client.Do("DEL", redisKeys...)
	return err
}

func (c *cache) StoreNil(cacheKey string) error {
	return c.Store(NewItemWithCustomTTL(cacheKey, []byte("null"), c.nilTTL))
}

func (c *cache) get(key string) (value any, err error) {
	client := c.connPool.Get()
	defer func() {
		_ = client.Close()
	}()

	err = client.Send("MULTI")
	if err != nil {
		return nil, err
	}
	err = client.Send("EXISTS", key)
	if err != nil {
		return nil, err
	}
	err = client.Send("GET", key)
	if err != nil {
		return nil, err
	}
	res, err := redigo.Values(client.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	val, ok := res[0].(int64)
	if !ok || val <= 0 {
		return nil, ErrKeyNotExist
	}

	return res[1], nil
}

func (c *cache) decideCacheTTL(item Item) (ttl int64) {
	if ttl = item.GetTTLInt64(); ttl > 0 {
		return
	}

	return int64(c.defaultTTL.Seconds())
}
