package lock

import (
	"github.com/go-redis/redis"
	rlock "github.com/bsm/redis-lock"
	"time"
	"st-go/errors"
)

var (
	ErrorLockFailed = errors.New("Lock failed")
)

type Options struct {
	Timeout   time.Duration
	RedisAddr string
	Prefix    string
}

type Locker interface {
	Lock(string) (Unlocker, error)
	LockWithTimeout(string, time.Duration) (Unlocker, error)
	Close() error
}

type Unlocker interface {
	Unlock() error
}

type unlocker struct {
	*rlock.Locker
}

func New(options Options) (Locker, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: options.RedisAddr,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	timeout := options.Timeout
	if timeout == 0 {
		timeout = time.Minute * 2
	}

	l := &lock{
		redisClient: redisClient,
		timeout:     timeout,
		prefix: options.Prefix,
	}
	return l, nil
}

type lock struct {
	redisClient *redis.Client
	timeout     time.Duration
	prefix      string
}

func (l *lock) Lock(key string) (Unlocker, error) {
	return l.LockWithTimeout(key, l.timeout)
}

func (l *lock) LockWithTimeout(key string, timeout time.Duration) (Unlocker, error) {
	key = l.prefix + key
	locker, err := rlock.Obtain(l.redisClient, key, &rlock.Options{LockTimeout: timeout})
	if err != nil {
		return nil, err
	}
	isLocked, err := locker.Lock()
	if err != nil {
		return nil, err
	}
	if !isLocked {
		return nil, ErrorLockFailed
	}
	return unlocker{locker}, nil
}

func (l *unlocker) Unlock() error {
	return l.Locker.Unlock()
}

func (l *lock) Close() error {
	return l.redisClient.Close()
}
