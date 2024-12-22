//go:build unit || integration

package common

import (
	"os"
	"sync"
)

var mockedEnv map[string]string

func init() {
	sync.OnceFunc(func() {
		mockedEnv = make(map[string]string)
	})()
}

func Putenv(key, value string) {
	mockedEnv[key] = value
}

func Delenv(key string) {
	delete(mockedEnv, key)
}

func Getenv(key string) string {
	if value, ok := mockedEnv[key]; ok {
		return value
	}
	return os.Getenv(key)
}
