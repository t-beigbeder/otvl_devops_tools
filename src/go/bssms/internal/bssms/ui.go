package bssms

import (
	"context"
	"net"
	"strconv"
)

type BaseConfig struct {
	Ctx context.Context
}

func (bc *BaseConfig) GetContext() context.Context {
	if bc.Ctx == nil {
		return context.Background()
	}
	return bc.Ctx
}

type ProvisionerConfig struct {
	BaseConfig
	UnsafeTls    bool
	ProxyAddress string
}

type InstallerConfig struct {
	BaseConfig
	UnsafeTls    bool
	ProxyAddress string
	Installable
}

type ProxyConfig struct {
	BaseConfig
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
