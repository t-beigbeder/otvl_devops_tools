package bssms

import (
	"github.com/urfave/cli/v2"
	"net"
	"strconv"
)

const (
	BssmsAlpn        = "x-otvl-bssms-v0.1"
	ProvisionerHello = "PrHello"
	MsgClose         = "Close"
)

type ProvisionerConfig struct {
	UnsafeTls    bool
	ProxyAddress string
}

func GetProvisionerConfig(cc *cli.Context) *ProvisionerConfig {
	return cc.App.Metadata["config"].(*ProvisionerConfig)
}

type InstallerConfig struct {
	UnsafeTls    bool
	ProxyAddress string
}

func GetInstallerConfig(cc *cli.Context) *InstallerConfig {
	return cc.App.Metadata["config"].(*InstallerConfig)
}

type ProxyConfig struct {
	UnsafeTls  bool
	ListenAddr string
	Host       string
}

func GetProxyConfig(cc *cli.Context) *ProxyConfig {
	return cc.App.Metadata["config"].(*ProxyConfig)
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
