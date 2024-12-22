package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"path/filepath"
)

type ProvisionerHost struct {
	PrivateKey string
}

type InstallHost struct {
	bssms.Installable `json:",inline" yaml:",inline"`
	PrivateKey        string `json:"privateKey,omitempty" yaml:"privateKey,omitempty"`
	PubKey            string `json:"pubKey,omitempty" yaml:"pubKey,omitempty"`
}

func ihfPath(optConfigDir string) (string, error) {
	r, err := bssms.GetConfigDir(optConfigDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(r, "installHosts.yaml"), nil
}

func LoadInstallHosts(optConfigDir string) ([]InstallHost, error) {
	p, err := ihfPath(optConfigDir)
	if err != nil {
		return nil, err
	}
	var ihs = []InstallHost{}
	if err := common.YamlLoad(p, &ihs); err != nil {
		return nil, err
	}
	return ihs, nil
}

func StoreInstallHosts(optConfigDir string, ihs []InstallHost) error {
	p, err := ihfPath(optConfigDir)
	if err != nil {
		return err
	}
	return common.YamlStore(p, ihs)
}
