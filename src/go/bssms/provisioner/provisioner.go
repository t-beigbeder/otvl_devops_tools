package provisioner

import (
	bssms "bssms/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/quic-go/quic-go"
)

func installing(sbr *bufio.Reader, cStream quic.Stream, iin bssms.Installable, ihs []InstallHost) {
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
	if err := common.WriteCommandToStream(sbr, cStream, bssms.ProvisionerInstall, pin, ""); err != nil {
		getLogger().Error("remote install failed", "err", err, "iin", iin)
		return
	}
	getLogger().Info("remote install started", "iin", iin)
	return
}

func installed(iin bssms.Installable, ihs []InstallHost) {
	ok, ih := matchFrom(iin, ihs)
	if !ok {
		getLogger().Error("installed does not match configured ones", "iin", iin)
		return
	}
	ih.Installed = true
	getLogger().Info("remote install achieved", "iin", iin)
}

func provision(ctx context.Context, conn quic.Connection, cStream quic.Stream, ihs []InstallHost) error {
	var (
		err            error
		rcmd           string
		allProvisioned bool
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
	for !allProvisioned {
		in := bssms.Installable{}
		rcmd, err = esr.ReadString('\n')
		if err != nil {
			return err
		}
		getLogger().Debug("reading proxy event", "rcmd", rcmd)
		if rcmd != bssms.ProvisionerEventInstallerUp && rcmd != bssms.ProvisionerEventInstalled {
			return fmt.Errorf("unexpected rcmd: %s", rcmd)
		}
		if err = common.ReadJsonFromStream(esr, &in); err != nil {
			getLogger().Error("ReadJsonFromStream", "err", err, "in", in)
			continue
		}
		getLogger().Debug("reading proxy event", "rcmd", rcmd, "in", in)
		if rcmd == bssms.ProvisionerEventInstallerUp {
			installing(sbr, cStream, in, ihs)
		} else {
			installed(in, ihs)
			allProvisioned = true
			for _, ih := range ihs {
				if !ih.Installed {
					getLogger().Debug("not all provisioned", "ih", ih)
					allProvisioned = false
				}
			}
		}
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
