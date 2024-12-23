package common

import (
	"fmt"
	"os"
	"testing"
)

func AsInterfaceSlice(values ...interface{}) []interface{} {
	return values
}

func SkipUnlessEnv(t *testing.T, key string) {
	if os.Getenv(key) == "" {
		t.Skip(fmt.Sprintf("skipped as env var %s not set", key))
	}
}
