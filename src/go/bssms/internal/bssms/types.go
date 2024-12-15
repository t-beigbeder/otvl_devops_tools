package bssms

import "github.com/urfave/cli/v2"

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
	UnsafeTls     bool
	ListenAddress string
}

func GetProxyConfig(cc *cli.Context) *ProxyConfig {
	return cc.App.Metadata["config"].(*ProxyConfig)
}
