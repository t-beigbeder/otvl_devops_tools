package provisioner

import (
	"bssms/bssms"
	"bssms/internal/common"
	"path/filepath"
	"sync"
)

type ProvisionerHost struct {
	PrivateKey string
}

type InstallHost struct {
	bssms.Installable `yaml:",inline"`
	PrivateKey        string            `yaml:"privateKey,omitempty"`
	PubKey            string            `yaml:"pubKey,omitempty"`
	Secrets           map[string]string `yaml:"secrets,omitempty"`
}

func matchFrom(iin bssms.Installable, pihs []*InstallHost) (bool, *InstallHost) {
	for i, ih := range pihs {
		if iin.Matches(ih.Installable) {
			return true, pihs[i]
		}
	}
	return false, nil
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

func ReadInstallHost(optConfigDir string, name string) (InstallHost, error) {
	ih := InstallHost{}
	p, err := ihfPath(optConfigDir)
	if err != nil {
		return ih, err
	}
	ihs, err := LoadYamlInstallHosts(p)
	if err != nil {
		return ih, err
	}
	for _, lih := range ihs {
		if lih.Name == name {
			return lih, nil
		}
	}
	return ih, nil
}

func SaveInstallHost(optConfigDir string, ih InstallHost) error {
	p, err := ihfPath(optConfigDir)
	if err != nil {
		return err
	}
	lockfile := sync.Mutex{}
	lockfile.Lock()
	defer lockfile.Unlock()
	ihs1, err := LoadYamlInstallHosts(p)
	ihs2 := []InstallHost{}
	found := false
	for _, lih := range ihs1 {
		if lih.Name == ih.Name {
			ihs2 = append(ihs2, ih)
			found = true
			continue
		}
		ihs2 = append(ihs2, lih)
	}
	if !found {
		ihs2 = append(ihs2, ih)
	}
	return common.YamlStore(p, ihs2)
}

func PInstallHosts(ihs []InstallHost) (pihs []*InstallHost) {
	for _, ih := range ihs {
		pihs = append(pihs, &ih)
	}
	return
}
