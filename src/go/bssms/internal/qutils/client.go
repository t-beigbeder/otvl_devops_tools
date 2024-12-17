package qutils

import (
	"bssms/internal/tlsutils"
	"context"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"time"
)

func GetQuicConn(addr string, alpn string) (quic.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := quic.DialAddr(ctx, addr,
		tlsutils.GetUnsafeTlsConfigClient(alpn), // TODO: configure TLS
		&quic.Config{Tracer: qlog.DefaultConnectionTracer},
	)
	return conn, err
}
