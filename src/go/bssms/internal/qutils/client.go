package qutils

import (
	"bssms/internal/tlsutils"
	"context"
	"github.com/quic-go/quic-go"
	"time"
)

func GetQuicConn(addr string) (*quic.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := quic.DialAddr(ctx, addr, tlsutils.GetUnsafeTlsConfigClient(), nil) // TODO: configure TLS
	return &conn, err
}