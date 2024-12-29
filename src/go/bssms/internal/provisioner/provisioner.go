package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"bufio"
	"fmt"
	"github.com/quic-go/quic-go"
)

func provision(stream quic.Stream, ihs []InstallHost) error {
	_, err := stream.Write([]byte(bssms.ProvisionerHello))
	if err != nil {
		return err
	}
	rs := bufio.NewReaderSize(stream, bssms.CtrlMsgMaxLn)
	cmd, err := rs.ReadString('\n')
	if err != nil {
		return err
	}
	if cmd != bssms.ProxyHello {
		return fmt.Errorf("invalid protocol command %s", cmd)
	}
	for _, ih := range ihs {
		bs, err := ih.Installable.AsJson()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("installable %s: %v", ih.Name, err))
		}
		_, err = stream.Write([]byte(bssms.ProvisionerInstallable))
		if err != nil {
			return err
		}
		_, err = stream.Write([]byte(fmt.Sprintf("%d\n", len(bs))))
		if err != nil {
			return err
		}
		_, err = stream.Write(bs)
		if err != nil {
			return err
		}
	}
	_, err = stream.Write([]byte(bssms.ApplicationClose))
	if err != nil {
		return err
	}
	return nil
}

func run(config *bssms.ProvisionerConfig, ihs []InstallHost) error {
	conn, err := qutils.GetQuicConn(config.ProxyAddress, bssms.BssmsAlpn)
	if err != nil {
		return err
	}
	defer conn.CloseWithError(0, "")
	stream, err := conn.OpenStreamSync(config.GetContext())
	if err != nil {
		return err
	}
	getLogger().Info("OpenStreamSync", "sid", stream.StreamID())
	defer stream.Close()
	return provision(stream, ihs)
}

func Run(config *bssms.ProvisionerConfig, optConfigDir string, ss []string) error {
	ihs, err := LoadFilteredInstallHosts(optConfigDir, ss)
	if err != nil {
		return err
	}
	return run(config, ihs)
}
