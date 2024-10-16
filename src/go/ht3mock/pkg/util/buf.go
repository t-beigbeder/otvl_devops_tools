package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
)

type randGenRC struct {
	size     int
	position int
	logger   *slog.Logger
}

const bufSize = 1024

func (r *randGenRC) Read(p []byte) (int, error) {
	if r.position >= r.size {
		return 0, io.EOF
	}
	gs := bufSize
	if r.position+bufSize > r.size {
		gs = r.size - r.position
	}
	if gs > len(p) {
		gs = len(p)
	}
	n, err := rand.Read(p[:gs])
	if n != gs || err != nil {
		r.position += n
		return n, fmt.Errorf("in hmClient.run %d %v", n, err)
	}
	r.position += n
	return n, nil
}

func (r *randGenRC) Close() error {
	return nil
}

func NewRandGenerator(size int, logger *slog.Logger) io.ReadCloser {
	return &randGenRC{size: size, logger: logger}
}
