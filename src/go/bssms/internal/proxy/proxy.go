package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"bssms/internal/tlsutils"
	"bufio"
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"os"
)

func handle(config *bssms.ProxyConfig, conn quic.Connection) error {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()
	rs := bufio.NewReaderSize(stream, bssms.CtrlMsgMaxLn)
	var (
		opened  bool
		closing bool
		cmd     string
	)
	for !closing && err == nil {
		cmd, err = rs.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Fprintf(os.Stderr, "handle %s", cmd)
		if !opened && (cmd == bssms.ProvisionerHello || cmd == bssms.InstallerHello) {
			opened = true
			_, err = stream.Write([]byte(bssms.ProxyHello + "\n"))
			if err != nil {
				break
			}
			continue
		}
		if opened && cmd == bssms.ApplicationClose {
			closing = true
			continue
		}
		err = fmt.Errorf("invalid protocol command %s", cmd)
	}
	return err
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
		conn, err := ln.Accept(context.Background())
		if err != nil {
			return err
		}
		go func(conn quic.Connection) {
			if err := handle(config, conn); err != nil {
				if ae, ok := err.(*quic.ApplicationError); !ok || ae.ErrorCode != 0 {
					fmt.Fprintf(os.Stderr, "connection error %v\n", err)
					conn.CloseWithError(1, fmt.Sprintf("connection error %v", err))
					return
				}
			} else {
				conn.CloseWithError(0, "")
			}
		}(conn)
		defer conn.CloseWithError(0, "")
	}
}
