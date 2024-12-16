package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"bssms/internal/tlsutils"
	"context"
)

func RunProxy(config *bssms.ProxyConfig) error {
	cert, err := tlsutils.SelfSigned(config.Host)
	if err != nil {
		return err
	}
	ln, err := qutils.GetQuicListener(config.ListenAddr, cert, bssms.BssmsAlpn)
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
