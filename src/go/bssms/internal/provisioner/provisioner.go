package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bufio"
	"errors"
	"fmt"
	"github.com/quic-go/quic-go"
)

func install(cStream quic.Stream, in bssms.Installable) error {
	return nil
}

func provision(cStream, eStream quic.Stream, ihs []InstallHost) error {
	_, err := cStream.Write([]byte(bssms.ProvisionerHello))
	if err != nil {
		return err
	}
	csr := bufio.NewReaderSize(cStream, bssms.CtrlMsgMaxLn)
	cmd, err := csr.ReadString('\n')
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
	err = common.WriteCommandToStream(cStream, bssms.ProvisionerInstallables, ihs)
	if err != nil {
		return err
	}
	esr := bufio.NewReaderSize(eStream, bssms.CtrlMsgMaxLn)
	for i := 0; i < len(ins); i++ {
		in := bssms.Installable{}
		if err = common.ReadJsonFromStream(esr, &in); err != nil {
			return err
		}
		install(cStream, in)
	}
	_, err = cStream.Write([]byte(bssms.ApplicationClose))
	if err != nil {
		return err
	}
	cmd, err = csr.ReadString('\n')
	if err != nil {
		var ae *quic.ApplicationError
		if !errors.As(err, &ae) || !ae.Remote || ae.ErrorCode != 0 {
			return err
		}
		return nil
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
	cStream, err := conn.OpenStreamSync(config.GetContext())
	if err != nil {
		return err
	}
	getLogger().Info("OpenStreamSync", "cSid", cStream.StreamID())
	defer cStream.Close()
	eStream, err := conn.OpenStreamSync(config.GetContext())
	if err != nil {
		return err
	}
	getLogger().Info("OpenStreamSync", "eSid", eStream.StreamID())
	defer eStream.Close()
	return provision(cStream, eStream, ihs)
}

func Run(config *bssms.ProvisionerConfig, optConfigDir string, ss []string) error {
	ihs, err := LoadFilteredInstallHosts(optConfigDir, ss)
	if err != nil {
		return err
	}
	return run(config, ihs)
}
