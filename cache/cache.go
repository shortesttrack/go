package cache

import (
	"github.com/go-redis/redis"
	"time"
	"st-go/errors"
)

type Cache interface {
	Set(string, interface{}) error
	SetWithExpiration(string, interface{}, time.Duration) error
	Delete(...string) error
	String(string) (string, error)
	Close()
}

type cache struct {
	*redis.Client

	expiration time.Duration
}

type Options struct {
	RedisAddr string
	Expiration time.Duration
}

func New(options Options) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: options.RedisAddr,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &cache{client, options.Expiration}, nil
}

func (c *cache) Close() {
	c.Client.Close()
}

func (c *cache) Set(key string, value interface{}) error {
	return c.SetWithExpiration(key, value, c.expiration)
}

func (c *cache) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(key, value, expiration).Err()
}

func (c *cache) Delete(keys ...string) error {
	return c.Client.Del(keys...).Err()
}

func (c *cache) String(key string) (string, error) {
	result, err := c.Client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.NotFound
		}
		return "", err
	}
	return result, nil
}