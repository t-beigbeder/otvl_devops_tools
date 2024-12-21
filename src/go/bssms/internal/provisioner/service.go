package provisioner

import (
	"bssms/internal/bssms"
	"gopkg.in/yaml.v3"
	"io"
	"os"
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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	y, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var ihs []InstallHost = []InstallHost{}
	err = yaml.Unmarshal(y, &ihs)
	if err != nil {
		return nil, err
	}
	return ihs, nil
}

func StoreInstallHosts(path string, ihs []InstallHost) error {
	y, err := yaml.Marshal(ihs)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(y)
	if err != nil {
		return err
	}
	return nil
}
