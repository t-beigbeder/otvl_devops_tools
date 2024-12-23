package integration

import (
	"bssms/internal/bssms"
	"bssms/internal/proxy"
	"golang.org/x/net/context"
	"testing"
	"time"
)

//func runTentative(cer bssms.ContextSetter, f func(bssms.ContextSetter) error) error {
//	f(context.Background())
//	return nil
//}

func RunProxy(t *testing.T) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	config := bssms.ProxyConfig{
		BaseConfig: bssms.BaseConfig{Ctx: ctx},
		UnsafeTls:  true,
		ListenAddr: ":9443",
		Host:       "localhost",
	}
	//runTentative(&config, proxy.RunProxy)
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
