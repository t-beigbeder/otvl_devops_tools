package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bssms/internal/tlsutils"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/quic-go/quic-go"
)

func handle(config *bssms.ProxyConfig, conn quic.Connection) error {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()
	getLogger().Info("AcceptStream", "sid", stream.StreamID())
	rs := bufio.NewReaderSize(stream, common.CtrlMsgMaxLn)
	var (
		opened  bool
		isPr    bool
		isIn    bool
		closing bool
		cmd     string
	)
	for !closing && err == nil {
		cmd, err = rs.ReadString('\n')
		if err != nil {
			break
		}
		getLogger().Info("handle", "cmd", cmd)
		if !opened && (cmd == bssms.ProvisionerHello || cmd == bssms.InstallerHello) {
			opened = true
			isPr = cmd == bssms.ProvisionerHello
			isIn = cmd == bssms.InstallerHello
			_, err = stream.Write([]byte(bssms.ProxyHello))
			if err != nil {
				break
			}
			continue
		}
		if opened && cmd == bssms.ApplicationBye {
			closing = true
			continue
		}
		if opened {
			if isPr {
				ins, err := handlePrCmd(stream, cmd)
				if err != nil {
					break // FIXME forever
				}
				_ = ins
			}
			if isIn {
				err = handleInCmd(config, stream, cmd)
				if err != nil {
					break
				}
			}
			continue
		}
		err = fmt.Errorf("invalid protocol command %s", cmd)
	}
	if err != nil {
		return err
	}
	_, err = stream.Write([]byte(bssms.ProxyBye))
	return err
}

func RunProxy(config *bssms.ProxyConfig) error {
	cert, err := tlsutils.SelfSigned(config.Host)
	if err != nil {
		return err
	}
	ln, err := qutils.GetQuicListener(config.ListenAddr, cert, bssms.BssmsAlpn, getLogger())
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept(config.GetContext())
		if err != nil {
			return err
		}
		go func(conn quic.Connection) {
			if err := handle(config, conn); err != nil {
				var ae *quic.ApplicationError
				if !errors.As(err, &ae) || ae.ErrorCode != 0 {
					getLogger().Error("connection error", "err", err)
					conn.CloseWithError(1, fmt.Sprintf("connection error %v", err))
					return
				}
			} else {
				conn.CloseWithError(0, "")
			}
		}(conn)
		defer conn.CloseWithError(0, "")
		getLogger().Info("Accept", "RemoteAddr", conn.RemoteAddr())
	}
}
