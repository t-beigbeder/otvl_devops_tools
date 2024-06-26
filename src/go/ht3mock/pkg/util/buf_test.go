package util

import (
	"io"
	"log/slog"
	"testing"
)

func TestNewRandGenerator(t *testing.T) {
	for size := 1024; size < 50000000; size *= 2 {
		rg := NewRandGenerator(size, slog.Default())
		i, err := io.Copy(io.Discard, rg)
		if err != nil {
			t.Error(err)
		}
		if i != int64(size) {
			t.Errorf("i %d", i)
		}
	}
}
