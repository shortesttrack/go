package lock

import (
	"testing"
	"time"
)

const (
	testKey = "test_key"
)

func TestLocker_Lock(t *testing.T) {
	locker, err := New(Options{
		Timeout: time.Second * 120,
		RedisAddr: "localhost:6379",
		Prefix: "test_",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer locker.Close()

	lock, err := locker.Lock(testKey)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 2)

	_, err = locker.Lock(testKey)
	if err == nil {
		t.Error(`lock should fail`)
	}

	err = lock.Unlock()
	if err != nil {
		t.Error(err)
	}

	lock, err = locker.Lock(testKey)
	if err != nil {
		t.Error(err)
	}
	defer lock.Unlock()
}