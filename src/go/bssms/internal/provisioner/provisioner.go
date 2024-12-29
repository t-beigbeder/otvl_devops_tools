package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
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
	ins := []bssms.Installable{}
	for _, ih := range ihs {
		ins = append(ins, ih.Installable)
	}
	err = common.WriteCommandToStream(stream, bssms.ProvisionerInstallables, ihs)
	if err != nil {
		return err
	}
	_, err = stream.Write([]byte(bssms.ApplicationClose))
	if err != nil {
		return err
	}
	cmd, err = rs.ReadString('\n')
	if err != nil {
		return err
	}
	if cmd != bssms.ProxyBye {
		return fmt.Errorf("invalid protocol command %s", cmd)
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
