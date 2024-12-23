package bssms

import (
	"golang.org/x/net/context"
	"net"
	"strconv"
)

type BaseConfig struct {
	Ctx context.Context
}

func (bc *BaseConfig) SetContext(ctx context.Context) {
	bc.Ctx = ctx
}

type ContextSetter interface {
	SetContext(ctx context.Context)
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
	IPAddress    string
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
