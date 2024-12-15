package qutils

import (
	"bssms/internal/bssms"
	"crypto/tls"
	"github.com/quic-go/quic-go"
	"net"
)

func GetQuicListener(addr string) (*quic.Listener, error) {
	ip, port, err := bssms.GetIPPort(addr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return nil, err
	}
	qc := quic.Config{}
	tc := tls.Config{
		Certificates:       nil,  // TODO: configure it
		InsecureSkipVerify: true, // TODO: configure it
	}
	return quic.Listen(udpConn, &tc, &qc)
}
