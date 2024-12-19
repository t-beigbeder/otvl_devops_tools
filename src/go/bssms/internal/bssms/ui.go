package bssms

import (
	"net"
	"strconv"
)

type ProvisionerConfig struct {
	UnsafeTls    bool
	ProxyAddress string
}

type InstallerConfig struct {
	UnsafeTls    bool
	ProxyAddress string
	Installable
}

type ProxyConfig struct {
	UnsafeTls  bool
	ListenAddr string
	Host       string
}

func GetIPPort(addr string) (net.IP, int, error) {
	is, ps, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, 0, err
	}
	i := net.ParseIP(is)
	p, err := strconv.ParseInt(ps, 10, 16)
	if err != nil {
		return nil, 0, err
	}
	return i, int(p), nil
}
