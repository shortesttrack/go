package lock

import (
	"github.com/go-redis/redis"
	rlock "github.com/bsm/redis-lock"
	"time"
	"st-go/errors"
	"sync"
	"github.com/satori/go.uuid"
)

var (
	defaultTimeout = time.Minute * 5
	ErrorLockFailed = errors.New("Lock failed")
)

type Options struct {
	Timeout   time.Duration
	RedisAddr string
	Prefix    string
}

type Locker interface {
	Lock(string) (Lock, error)
	LockWithTimeout(string, time.Duration) (Lock, error)
	Close() error
}

type Lock interface {
	Unlock() error
}

type lock struct {
	*rlock.Locker
	*locker

	uuid uuid.UUID
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
		timeout = defaultTimeout
	}

	l := &locker{
		redisClient: redisClient,
		timeout:     timeout,
		prefix: options.Prefix,
		activeLocks: make(map[uuid.UUID]*lock),
	}
	return l, nil
}

type locker struct {
	redisClient *redis.Client
	timeout     time.Duration
	prefix      string

	mu sync.Mutex
	activeLocks map[uuid.UUID]*lock
}

func (l *locker) Lock(key string) (Lock, error) {
	return l.LockWithTimeout(key, l.timeout)
}

func (l *locker) LockWithTimeout(key string, timeout time.Duration) (Lock, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
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
	lock := &lock{locker, l, uuid}
	l.mu.Lock()
	l.activeLocks[lock.uuid] = lock
	l.mu.Unlock()
	return lock, nil
}

func (l *locker) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, lock := range l.activeLocks {
		lock.Locker.Unlock()
	}
	return l.redisClient.Close()
}

func (l *lock) Unlock() error {
	l.mu.Lock()
	delete(l.activeLocks, l.uuid)
	l.mu.Unlock()
	return l.Locker.Unlock()
}
