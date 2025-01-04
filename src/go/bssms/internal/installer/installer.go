package installer

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bssms/internal/qutils"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/quic-go/quic-go"
)

func handleEvents(config *bssms.InstallerConfig, conn quic.Connection) (map[string]string, error) {
	var (
		err        error
		rcmd       string
		encSecrets []byte
		jsSecrets  string
		secrets    map[string]string
	)

	eStream, err := conn.AcceptStream(config.GetContext())
	if err != nil {
		return nil, err
	}
	defer eStream.Close()
	getLogger().Info("AcceptStream", "eSid", eStream.StreamID())
	esr := bufio.NewReaderSize(eStream, common.CtrlDataMaxLn)
	rcmd, err = esr.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if rcmd != bssms.InstallerEventInstall {
		return nil, fmt.Errorf("unexpected rcmd: %s", rcmd)
	}
	if encSecrets, err = common.ReadBytesFromStream(esr); err != nil {
		return nil, err
	}
	getLogger().Debug("install", "encSecrets", encSecrets)
	if jsSecrets, err = common.DecryptMsg(encSecrets, config.PrivateKey); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(jsSecrets), &secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}

func output(config *bssms.InstallerConfig, secrets map[string]string) error {
	if config.JsonSecf != "" {
		return common.JsonStore(config.JsonSecf, secrets)
	} else if config.YamlSecf != "" {
		return common.YamlStore(config.YamlSecf, secrets)
	} else {
		for k, v := range secrets {
			fmt.Println(k, v)
		}
	}
	return nil
}

func install(config *bssms.InstallerConfig, conn quic.Connection, cStream quic.Stream) error {
	var (
		err     error
		rcmd    string
		secrets map[string]string
	)
	sbr := bufio.NewReaderSize(cStream, common.CtrlDataMaxLn)
	if err = common.WriteCommandToStream(sbr, cStream, bssms.InstallerHello, nil, bssms.ProxyHello); err != nil {
		return err
	}
	in := bssms.Installable{
		ServerUuid: config.ServerUuid,
		MacAddress: config.MacAddress,
		IPAddress:  config.IPAddress,
	}
	if err = common.WriteCommandToStream(sbr, cStream, bssms.InstallerInstallable, in, ""); err != nil {
		return err
	}

	secrets, err = handleEvents(config, conn)
	if err != nil {
		return err
	}
	getLogger().Debug("install", "secrets", secrets)
	if err = output(config, secrets); err != nil {
		return err
	}

	if err = common.WriteCommandToStream(sbr, cStream, bssms.InstallerInstalled, in, ""); err != nil {
		return err
	}

	rcmd, err = common.WriteByeCommandToStream(sbr, cStream, bssms.ApplicationBye, bssms.ProxyBye)
	if err != nil {
		return err
	}
	if rcmd == "" {
		getLogger().Info("install: remote closed connection")
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
