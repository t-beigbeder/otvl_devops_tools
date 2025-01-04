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
	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

type clCtxBaseType struct {
	cid         string
	lner        *listener
	established bool
	sbr         *bufio.Reader
}
type prCtxType struct {
	clCtxBaseType
	ins []bssms.Installable
	in  bssms.Installable
}

type inCtxType struct {
	clCtxBaseType
	in bssms.Installable
}

func handle(config *bssms.ProxyConfig, conn quic.Connection, lner *listener) error {
	cid := uuid.New().String()
	var (
		cStream quic.Stream
		prCtx   *prCtxType
		inCtx   *inCtxType
		err     error
		opened  bool
		closing bool
		cmd     string
	)
	cStream, err = conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	defer cStream.Close()
	getLogger().Info("AcceptStream", "sid", cStream.StreamID())
	sbr := bufio.NewReaderSize(cStream, common.CtrlDataMaxLn)
	for !closing && err == nil {
		cmd, err = sbr.ReadString('\n')
		if err != nil {
			break
		}
		getLogger().Info("handle", "cmd", cmd)
		if !opened && (cmd == bssms.ProvisionerHello || cmd == bssms.InstallerHello) {
			opened = true
			if cmd == bssms.ProvisionerHello {
				prCtx = &prCtxType{
					clCtxBaseType: clCtxBaseType{cid: cid, lner: lner, sbr: sbr},
					ins:           nil,
					in:            bssms.Installable{},
				}
			} else {
				inCtx = &inCtxType{
					clCtxBaseType: clCtxBaseType{cid: cid, lner: lner, sbr: sbr},
					in:            bssms.Installable{},
				}
			}
			_, err = cStream.Write([]byte(bssms.ProxyHello))
			if err != nil {
				break
			}
			err = lner.addConnection(cid, conn, prCtx != nil)
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
			if prCtx != nil {
				if err = handlePrCmds(prCtx, cmd); err != nil {
					break
				}
			} else {
				if err = handleInCmds(inCtx, cmd); err != nil {
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
	_, err = cStream.Write([]byte(bssms.ProxyBye))
	return err
}

func RunProxy(config *bssms.ProxyConfig) error {
	cert, err := tlsutils.SelfSigned(config.Host)
	if err != nil {
		return err
	}
	qln, err := qutils.GetQuicListener(config.ListenAddr, cert, bssms.BssmsAlpn, getLogger())
	if err != nil {
		return err
	}
	lner, err := makeListener(config.GetContext())
	if err != nil {
		return err
	}
	defer lner.close()
	for {
		conn, err := qln.Accept(config.GetContext())
		if err != nil {
			getLogger().Error("accept", "err", err)
			return err
		}
		go func(conn quic.Connection) {
			if err := handle(config, conn, lner); err != nil {
				var ae *quic.ApplicationError
				if !errors.As(err, &ae) || ae.ErrorCode != 0 {
					getLogger().Error("connection error", "err", err)
					conn.CloseWithError(1, fmt.Sprintf("connection error %v", err))
					return
				}
				getLogger().Info("application error", "err", err)
			} else {
				getLogger().Info("close connection without error", "conn", conn.RemoteAddr())
				if err = conn.CloseWithError(0, ""); err != nil {
					getLogger().Error("CloseWithError 0", "err", err)
				}
			}
		}(conn)
		defer func() {
			getLogger().Info("deferred close connection without error", "conn", conn.RemoteAddr())
			if err = conn.CloseWithError(0, ""); err != nil {
				getLogger().Error("deferred CloseWithError 0", "err", err)
			}
		}()
		getLogger().Info("Accept", "RemoteAddr", conn.RemoteAddr())
	}
}
