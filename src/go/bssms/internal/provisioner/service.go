package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
)

type ProvisionerHost struct {
	PrivateKey string
}

type InstallHost struct {
	bssms.Installable `json:",inline" yaml:",inline"`
	PrivateKey        string `json:"privateKey,omitempty" yaml:"privateKey,omitempty"`
	PubKey            string `json:"pubKey,omitempty" yaml:"pubKey,omitempty"`
}

func LoadInstallHosts(path string) ([]InstallHost, error) {
	var ihs = []InstallHost{}
	if err := common.YamlLoad(path, &ihs); err != nil {
		return nil, err
	}
	return ihs, nil
}

func StoreInstallHosts(path string, ihs []InstallHost) error {
	return common.YamlStore(path, ihs)
}
