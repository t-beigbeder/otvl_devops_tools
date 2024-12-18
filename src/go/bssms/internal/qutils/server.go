package qutils

import (
	"bssms/internal/bssms"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"net"
	"os"
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
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			fmt.Fprintf(os.Stderr, "connection from client %+v\n", info)
			return nil, nil
		},
	}
	return quic.Listen(udpConn, &tc, &qc)
}
