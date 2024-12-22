//go:build !(unit || integration)

package common

import "os"

func Getenv(key string) string {
	return os.Getenv(key)
}
