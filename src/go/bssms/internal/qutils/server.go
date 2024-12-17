package qutils

import (
	"bssms/internal/bssms"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"net"
)

func GetQuicListener(addr string, cert *tls.Certificate, alpn string) (*quic.Listener, error) {
	ip, port, err := bssms.GetIPPort(addr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return nil, err
	}
	qc := quic.Config{Tracer: qlog.DefaultConnectionTracer}
	tc := tls.Config{
		Certificates: []tls.Certificate{*cert},
		NextProtos:   []string{alpn},
	}
	return quic.Listen(udpConn, &tc, &qc)
}
