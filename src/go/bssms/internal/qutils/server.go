package qutils

import (
	"context"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"net"
)

func RunServer() error {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 1234})
	if err != nil {
		return err
	}
	qc := quic.Config{}
	tc := tls.Config{
		Certificates: nil,
	}
	ln, err := quic.Listen(udpConn, &tc, &qc)
	for {
		conn, err := ln.Accept(context.TODO())
		if err != nil {
			return err
		}
		_ = conn
	}
}
