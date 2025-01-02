package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bufio"
	"context"
	"fmt"
	"github.com/quic-go/quic-go"
	"sync"
)

type connd struct {
	conn   quic.Connection
	stream quic.Stream
	sbr    *bufio.Reader
	isPr   bool
	ins    []bssms.Installable
	isIn   bool
	in     bssms.Installable
}

type listener struct {
	ctx    context.Context
	done   chan struct{}
	mux    sync.Mutex
	connds map[string]*connd
}

func (lner *listener) listen() {
	getLogger().Info("start listening on listener")
	for {
		select {
		case <-lner.done:
			getLogger().Info("listener closed")
			close(lner.done)
			return
		}
	}
}

func (lner *listener) close() {
	getLogger().Info("closing listener")
	lner.done <- struct{}{}
}

func (lner *listener) addConnection(cid string, conn quic.Connection, isPr, isIn bool) error {
	lner.mux.Lock()
	defer lner.mux.Unlock()
	stream, err := conn.OpenStreamSync(lner.ctx)
	if err != nil {
		return err
	}
	getLogger().Info("listener.OpenStreamSync", "cid", cid, "sid", stream.StreamID(), "isPr", isPr, "isIn", isIn)
	sbr := bufio.NewReaderSize(stream, common.CtrlDataMaxLn)
	lner.connds[cid] = &connd{
		conn:   conn,
		stream: stream,
		sbr:    sbr,
		isPr:   isPr,
		isIn:   isIn,
	}
	return nil
}

func (lner *listener) checkAndSendPrEvInUp(pcd *connd, pin, iin bssms.Installable) {
	if pin.ServerUuid == iin.ServerUuid &&
		pin.MacAddress == iin.MacAddress &&
		(pin.IPIntAddress == iin.IPAddress || pin.IPExtAddress == iin.IPAddress) {
		if err := common.WriteCommandToStream(pcd.sbr, pcd.stream, bssms.ProvisionerEventInstallerUp, pin, ""); err != nil {
			getLogger().Error("fail to send "+bssms.ProvisionerEventInstallerUp, "pin", pin, "err", err)
			return
		}
		getLogger().Info("sent "+bssms.ProvisionerEventInstallerUp, "pin", pin)
		return
	}
}

func (lner *listener) checkInstallerUp(pcd *connd) {
	for _, icd := range lner.connds {
		if icd.isPr {
			continue
		}
		for _, pin := range pcd.ins {
			lner.checkAndSendPrEvInUp(pcd, pin, icd.in)
		}
	}
}

func (lner *listener) checkProvisionerUp(icd *connd) {
	for _, pcd := range lner.connds {
		if pcd.isIn {
			continue
		}
		lner.checkAndSendPrEvInUp(pcd, pcd.in, icd.in)
	}
}

func (lner *listener) provisionerReadyEvent(cid string, ins []bssms.Installable) error {
	lner.mux.Lock()
	defer lner.mux.Unlock()
	pcd, ok := lner.connds[cid]
	if !ok {
		return fmt.Errorf("no connection found for cid: %s", cid)
	}
	pcd.ins = ins
	lner.checkInstallerUp(pcd)
	return nil
}

func (lner *listener) installerReadyEvent(cid string, in bssms.Installable) error {
	lner.mux.Lock()
	defer lner.mux.Unlock()
	icd, ok := lner.connds[cid]
	if !ok {
		return fmt.Errorf("no connection found for cid: %s", cid)
	}
	icd.in = in
	lner.checkProvisionerUp(icd)
	return nil
}

func makeListener(ctx context.Context) (*listener, error) {
	lner := &listener{
		ctx:    ctx,
		done:   make(chan struct{}),
		connds: make(map[string]*connd),
	}
	go lner.listen()
	return lner, nil
}
