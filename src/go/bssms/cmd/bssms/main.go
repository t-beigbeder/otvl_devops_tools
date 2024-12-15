package main

import (
	"bssms/internal/bssms"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"os"
)

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
					bssms.GetProvisionerConfig(cc).ProxyAddress = hp
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ut",
				Usage: "UnsafeTls",
				Action: func(cc *cli.Context, b bool) error {
					bssms.GetProvisionerConfig(cc).UnsafeTls = b
					return nil
				},
			},
		},
		Action: func(cc *cli.Context) error {
			config := bssms.GetProvisionerConfig(cc)
			_ = config
			return nil
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
					bssms.GetInstallerConfig(cc).ProxyAddress = hp
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ut",
				Usage: "UnsafeTls",
				Action: func(cc *cli.Context, b bool) error {
					bssms.GetInstallerConfig(cc).UnsafeTls = b
					return nil
				},
			},
		},
		Action: func(cc *cli.Context) error {
			config := bssms.GetInstallerConfig(cc)
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
				Name:     "pxa",
				Required: true,
				Usage:    "listen address: 'host:port' or 'ip:port'",
				Action: func(cc *cli.Context, hp string) error {
					if _, _, err := net.SplitHostPort(hp); err != nil {
						return err
					}
					bssms.GetProxyConfig(cc).ListenAddress = hp
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "ut",
				Usage: "UnsafeTls",
				Action: func(cc *cli.Context, b bool) error {
					bssms.GetProxyConfig(cc).UnsafeTls = b
					return nil
				},
			},
		},
		Action: func(cc *cli.Context) error {
			config := bssms.GetProxyConfig(cc)
			_ = config
			return nil
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
