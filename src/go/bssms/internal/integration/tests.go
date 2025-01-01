package integration

import (
	"bssms/internal/bssms"
	"bssms/internal/proxy"
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	ProxyHost = "localhost"
	ProxyPort = "9443"
)

func ProxyAddress() string { return fmt.Sprintf("%s:%s", ProxyHost, ProxyPort) }

func RunProxy(t *testing.T) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	config := bssms.ProxyConfig{
		BaseConfig: bssms.BaseConfig{Ctx: ctx},
		UnsafeTls:  true,
		ListenAddr: fmt.Sprintf(":%s", ProxyPort),
		Host:       ProxyHost,
	}
	var bgErr error
	go func() {
		err := proxy.RunProxy(&config)
		if err != nil {
			bgErr = err
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if bgErr != nil {
		cancel()
		return nil, bgErr
	}
	return cancel, nil
}

func GetTestDataDir() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("path discovery")
	}
	return filepath.Join(filepath.Dir(thisFile), "testdata")
}
