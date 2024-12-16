package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"bssms/internal/tlsutils"
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"os"
)

func handle(config *bssms.ProxyConfig, cnc quic.Connection) error {
	stream, err := cnc.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()
	buf := make([]byte, 128)
	_, err = stream.Read(buf)
	if err != nil {
		return err
	}
	_, err = stream.Write(buf)
	if err != nil {
		return err
	}
	_, err = stream.Read(buf)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "handle stop\n")
	return nil
}

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
		cnc, err := ln.Accept(context.TODO())
		if err != nil {
			return err
		}
		go func(cnc quic.Connection) {
			if err := handle(config, cnc); err != nil {
				fmt.Fprintf(os.Stderr, "connection error %v\n", err)
				cnc.CloseWithError(1, fmt.Sprintf("connection error %v", err))
			} else {
				cnc.CloseWithError(0, "")
			}
		}(cnc)
		defer cnc.CloseWithError(0, "")
	}
}
