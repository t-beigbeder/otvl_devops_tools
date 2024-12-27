package integration

import (
	"bssms/internal/bssms"
	"bssms/internal/proxy"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func RunProxy(t *testing.T) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	config := bssms.ProxyConfig{
		BaseConfig: bssms.BaseConfig{Ctx: ctx},
		UnsafeTls:  true,
		ListenAddr: ":9443",
		Host:       "localhost",
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
