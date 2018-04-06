package util

import (
	"testing"
	"os"
	"st-go/errors"
)

func TestGetEnv(t *testing.T) {
	testKey := "test-key"
	testValue := "test-value"
	fallbackValue := "fallback-value"
	err := os.Setenv(testKey, testValue)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(testKey)

	value := GetEnv(testKey, fallbackValue)
	if value != testValue {
		t.Error(errors.New("bad value"))
	}

}