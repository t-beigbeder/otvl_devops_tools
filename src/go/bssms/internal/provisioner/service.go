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
	bssms.Installable `yaml:",inline"`
	PrivateKey        string `yaml:"privateKey,omitempty"`
	PubKey            string `yaml:"pubKey,omitempty"`
}

func ihfPath(optConfigDir string) (string, error) {
	r, err := bssms.GetConfigDir(optConfigDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(r, "installHosts.yaml"), nil
}

func LoadYamlInstallHosts(path string) ([]InstallHost, error) {
	var ihs = []InstallHost{}
	if err := common.YamlLoad(path, &ihs); err != nil {
		return nil, err
	}
	return ihs, nil
}

func LoadInstallHosts(optConfigDir string) ([]InstallHost, error) {
	p, err := ihfPath(optConfigDir)
	if err != nil {
		return nil, err
	}
	return LoadYamlInstallHosts(p)
}

func LoadFilteredInstallHosts(optConfigDir string, fhn []string) ([]InstallHost, error) {
	ihs, err := LoadInstallHosts(optConfigDir)
	if err != nil {
		return nil, err
	}
	if len(fhn) == 0 {
		return ihs, err
	}
	fhs := []InstallHost{}
	for _, ih := range ihs {
		for _, fh := range fhn {
			if fh == ih.Name {
				fhs = append(fhs, ih)
				break
			}
		}
	}
	return fhs, err
}

func StoreInstallHosts(optConfigDir string, ihs []InstallHost) error {
	p, err := ihfPath(optConfigDir)
	if err != nil {
		return err
	}
	return common.YamlStore(p, ihs)
}
