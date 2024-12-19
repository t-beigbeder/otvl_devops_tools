package main

import (
	"bssms/internal/bssms"
	"bssms/internal/provisioner"
	"bssms/internal/proxy"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"os"
)

func getProvisionerConfig(cc *cli.Context) *bssms.ProvisionerConfig {
	return cc.App.Metadata["config"].(*bssms.ProvisionerConfig)
}

func getInstallerConfig(cc *cli.Context) *bssms.InstallerConfig {
	return cc.App.Metadata["config"].(*bssms.InstallerConfig)
}

func GetProxyConfig(cc *cli.Context) *bssms.ProxyConfig {
	return cc.App.Metadata["config"].(*bssms.ProxyConfig)
}

func getPrCmd() *cli.Command {
	return &cli.Command{
		Name:        "pr",
		Description: "provisioner",
		Before: func(cc *cli.Context) error {
			cc.App.Metadata["config"] = &bssms.ProvisionerConfig{}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "pxa",
				Required: true,
				Usage:    "proxy address: 'host:port' or 'ip:port'",
				Action: func(cc *cli.Context, hp string) error {
					if _, _, err := net.SplitHostPort(hp); err != nil {
						return err
					}
					getProvisionerConfig(cc).ProxyAddress = hp
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ut",
				Usage: "UnsafeTls",
				Action: func(cc *cli.Context, b bool) error {
					getProvisionerConfig(cc).UnsafeTls = b
					return nil
				},
			},
		},
		Action: func(cc *cli.Context) error {
			config := getProvisionerConfig(cc)
			err := provisioner.Run(config)
			return err
		},
	}
}

func getInCmd() *cli.Command {
	return &cli.Command{
		Name:        "in",
		Description: "installer",
		Before: func(cc *cli.Context) error {
			cc.App.Metadata["config"] = &bssms.InstallerConfig{}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "pxa",
				Required: true,
				Usage:    "proxy address: 'host:port' or 'ip:port'",
				Action: func(cc *cli.Context, hp string) error {
					if _, _, err := net.SplitHostPort(hp); err != nil {
						return err
					}
					getInstallerConfig(cc).ProxyAddress = hp
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "uuid",
				Required: true,
				Usage:    "host server uuid",
				Action: func(cc *cli.Context, s string) error {
					getInstallerConfig(cc).ServerUuid = s
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "mac",
				Required: true,
				Usage:    "host mac address",
				Action: func(cc *cli.Context, s string) error {
					getInstallerConfig(cc).MacAddress = s
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "ip",
				Required: true,
				Usage:    "host ip address used with proxy",
				Action: func(cc *cli.Context, s string) error {
					getInstallerConfig(cc).IPAddress = s
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ut",
				Usage: "UnsafeTls",
				Action: func(cc *cli.Context, b bool) error {
					getInstallerConfig(cc).UnsafeTls = b
					return nil
				},
			},
		},
		Action: func(cc *cli.Context) error {
			config := getInstallerConfig(cc)
			_ = config
			return nil
		},
	}
}

func getPxCmd() *cli.Command {
	return &cli.Command{
		Name:        "px",
		Description: "proxy",
		Before: func(cc *cli.Context) error {
			cc.App.Metadata["config"] = &bssms.ProxyConfig{}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "la",
				Required: true,
				Usage:    "listen address: ':port' or 'ip:port'",
				Action: func(cc *cli.Context, addr string) error {
					if _, _, err := bssms.GetIPPort(addr); err != nil {
						return err
					}
					GetProxyConfig(cc).ListenAddr = addr
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "host",
				Required: true,
				Usage:    "host for the certificate",
				Action: func(cc *cli.Context, host string) error {
					GetProxyConfig(cc).Host = host
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ut",
				Usage: "UnsafeTls",
				Action: func(cc *cli.Context, b bool) error {
					GetProxyConfig(cc).UnsafeTls = b
					return nil
				},
			},
		},
		Action: func(cc *cli.Context) error {
			config := GetProxyConfig(cc)
			return proxy.RunProxy(config)
		},
	}
}

func main() {
	app := &cli.App{
		Name:  "bssms",
		Usage: "use one subcommand",
		Commands: []*cli.Command{
			getPrCmd(),
			getInCmd(),
			getPxCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
