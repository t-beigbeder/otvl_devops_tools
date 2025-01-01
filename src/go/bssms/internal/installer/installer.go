package installer

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bufio"
	"github.com/quic-go/quic-go"
	"golang.org/x/net/context"
)

func install(config *bssms.InstallerConfig, conn quic.Connection, cStream quic.Stream) error {
	var (
		err  error
		rcmd string
	)
	sbr := bufio.NewReaderSize(cStream, common.CtrlDataMaxLn)
	if err = common.WriteCommandToStream(sbr, cStream, bssms.InstallerHello, nil, bssms.ProxyHello); err != nil {
		return err
	}
	in := bssms.Installable{
		Name:       "",
		ServerUuid: "",
		MacAddress: "",
	}
	err = common.WriteCommandToStream(sbr, cStream, bssms.ProvisionerInstallables, in, "")
	if err != nil {
		return err
	}
	eStream, err := conn.AcceptStream(config.GetContext())
	if err != nil {
		return err
	}
	defer eStream.Close()
	getLogger().Info("AcceptStream", "eSid", cStream.StreamID())
	esr := bufio.NewReaderSize(eStream, common.CtrlMsgMaxLn)
	for i := 0; i < len(ins); i++ {
		in := bssms.Installable{}
		if err = common.ReadJsonFromStream(esr, &in); err != nil {
			return err
		}
		install(cStream, in)
	}

	rcmd, err = common.WriteByeCommandToStream(sbr, cStream, bssms.ApplicationBye, bssms.ProxyBye)
	if err != nil {
		return err
	}
	if rcmd == "" {
		getLogger().Info("provision: remote closed connection")
	}
	return nil
}

func Run(config *bssms.InstallerConfig) error {
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
	return install(config, conn, cStream)
}
