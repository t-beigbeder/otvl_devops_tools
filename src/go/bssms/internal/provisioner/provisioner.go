package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bufio"
	"github.com/quic-go/quic-go"
)

func install(cStream quic.Stream, in bssms.Installable) error {
	return nil
}

func provision(cStream, eStream quic.Stream, ihs []InstallHost) error {
	var (
		err  error
		rcmd string
	)
	sbr := bufio.NewReaderSize(cStream, common.CtrlDataMaxLn)
	if err = common.WriteCommandToStream(sbr, cStream, bssms.ProvisionerHello, nil, bssms.ProxyHello); err != nil {
		return err
	}
	ins := []bssms.Installable{}
	for _, ih := range ihs {
		ins = append(ins, ih.Installable)
	}
	err = common.WriteCommandToStream(sbr, cStream, bssms.ProvisionerInstallables, ihs, "")
	if err != nil {
		return err
	}
	// esr := bufio.NewReaderSize(eStream, common.CtrlMsgMaxLn)
	//for i := 0; i < len(ins); i++ {
	//	in := bssms.Installable{}
	//	if err = common.ReadJsonFromStream(esr, &in); err != nil {
	//		return err
	//	}
	//	install(cStream, in)
	//}

	rcmd, err = common.WriteByeCommandToStream(sbr, cStream, bssms.ApplicationBye, bssms.ProxyBye)
	if err != nil {
		return err
	}
	if rcmd == "" {
		getLogger().Info("provision: remote closed connection")
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
