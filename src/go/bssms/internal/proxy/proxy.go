package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"context"
)

func RunProxy(config *bssms.ProxyConfig) error {
	ln, err := qutils.GetQuicListener(config.ListenAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept(context.TODO())
		if err != nil {
			return err
		}
		_ = conn
	}
}
