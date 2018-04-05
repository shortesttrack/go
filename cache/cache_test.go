package cache

import (
	"testing"
	"time"
	"st-go/errors"
)

const (
	testKey = "test-key"
	testValue = "test-value"
	testExpiration = time.Second * 5
)

func makeTestCache() (Cache, error) {
	return New(Options{
		RedisAddr: "localhost:6379",
		Expiration: time.Second * 120,
	})
}

func TestCache_Set(t *testing.T) {
	cache, err := makeTestCache()
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()

	err = cache.Set(testKey, testValue)
	if err != nil {
		t.Error(err)
	}

	value, err := cache.String(testKey)
	if err != nil {
		t.Error(err)
	}
	if value != testValue {
		t.Error(errors.New("result value not matching original"))
	}
}

func TestCache_SetWithExpiration(t *testing.T) {
	cache, err := makeTestCache()
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()

	err = cache.SetWithExpiration(testKey, testValue, testExpiration)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(testExpiration / 2)

	value, err := cache.String(testKey)
	if err != nil {
		t.Error(err)
	}

	if value != testValue {
		t.Error(errors.New("result value not matching original"))
	}

	time.Sleep(testExpiration / 2 + 1)

	value, err = cache.String(testKey)
	if err == nil {
		t.Error(errors.New("should be not found"))
	}
	if err != errors.NotFound {
		t.Error(err)
	}

}

func TestCache_Delete(t *testing.T) {
	cache, err := makeTestCache()
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()

	err = cache.Set(testKey, testValue)
	if err != nil {
		t.Fatal(err)
	}

	err = cache.Delete(testKey)
	if err != nil {
		t.Error(err)
	}

	_, err = cache.String(testKey)
	if err == nil {
		t.Error(errors.New("should be not found"))
	}
	if err != errors.NotFound {
		t.Error(err)
	}
}