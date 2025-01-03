package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/quic-go/quic-go"
)

func install(sbr *bufio.Reader, cStream quic.Stream, iin bssms.Installable, ihs []InstallHost) {
	ok, ih := matchFrom(iin, ihs)
	if !ok {
		getLogger().Error("installable does not match configured ones", "iin", iin)
		return
	}
	js, err := json.Marshal(ih.Secrets)
	if err != nil {
		getLogger().Error("json encoding failed", "err", err, "iin", iin)
		return
	}
	ebs, err := common.EncryptMsg(string(js), ih.PubKey)
	if err != nil {
		getLogger().Error("encryption failed", "err", err, "iin", iin)
		return
	}
	pin := iin
	pin.EncSecrets = string(ebs)
	if err := common.WriteCommandToStream(sbr, cStream, bssms.ProvisionerInstall, pin, bssms.ProxyHostInstalled); err != nil {
		getLogger().Error("remote install failed", "err", err, "iin", iin)
		return
	}
	getLogger().Info("remote install succeeded", "iin", iin)
	return
}

func provision(ctx context.Context, conn quic.Connection, cStream quic.Stream, ihs []InstallHost) error {
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
	eStream, err := conn.AcceptStream(ctx)
	if err != nil {
		return err
	}
	defer eStream.Close()
	getLogger().Info("AcceptStream", "eSid", cStream.StreamID())
	esr := bufio.NewReaderSize(eStream, common.CtrlDataMaxLn)
	for i := 0; i < len(ins); i++ {
		in := bssms.Installable{}
		rcmd, err = esr.ReadString('\n')
		if err != nil {
			return err
		}
		if rcmd != bssms.ProvisionerEventInstallerUp {
			return fmt.Errorf("unexpected rcmd: %s", rcmd)
		}
		if err = common.ReadJsonFromStream(esr, &in); err != nil {
			getLogger().Error("ReadJsonFromStream", "err", err, "in", ins[i])
			continue
		}
		install(sbr, cStream, in, ihs)
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

func RunIhs(config *bssms.ProvisionerConfig, ihs []InstallHost) error {
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
	return provision(config.GetContext(), conn, cStream, ihs)
}

func Run(config *bssms.ProvisionerConfig, optConfigDir string, ss []string) error {
	ihs, err := LoadFilteredInstallHosts(optConfigDir, ss)
	if err != nil {
		return err
	}
	return RunIhs(config, ihs)
}
