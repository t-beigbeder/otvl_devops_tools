package proxy

import (
	"bssms/internal/common"
	"bufio"
	"github.com/quic-go/quic-go"
	"golang.org/x/net/context"
	"sync"
)

type connd struct {
	conn   quic.Connection
	stream quic.Stream
	sbr    *bufio.Reader
	isPr   bool
	isIn   bool
}

type listener struct {
	ctx   context.Context
	done  chan struct{}
	mux   sync.Mutex
	conns map[string]connd
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
	getLogger().Info("listener.OpenStreamSync", "cid", cid, "sid", stream.StreamID())
	sbr := bufio.NewReaderSize(stream, common.CtrlDataMaxLn)
	lner.conns[cid] = connd{
		conn:   conn,
		stream: stream,
		sbr:    sbr,
		isPr:   isPr,
		isIn:   isIn,
	}
	return nil
}

func makeListener(ctx context.Context) (*listener, error) {
	lner := &listener{
		ctx:   ctx,
		done:  make(chan struct{}),
		conns: make(map[string]connd),
	}
	go lner.listen()
	return lner, nil
}
