package provisioner

import (
	"bssms/bssms"
	"bssms/internal/common"
	"sync"
)

func newIh(name string) (ih InstallHost, err error) {
	var (
		pub, pri string
	)
	if pub, pri, err = common.NewKeyPair(); err != nil {
		return
	}
	ih = InstallHost{
		Installable: bssms.Installable{Name: name},
		PubKey:      pub,
		PrivateKey:  pri,
	}
	return
}

func RunPhase0(optConfigDir string, names []string) error {
	var ihs []InstallHost
	for _, name := range names {
		ih, err := newIh(name)
		if err != nil {
			return err
		}
		ihs = append(ihs, ih)
	}
	return StoreInstallHosts(optConfigDir, ihs)
}

var fileLock sync.Mutex

func MergePhase0(optConfigDir string, name string) error {
	fileLock.Lock()
	defer fileLock.Unlock()
	ihs, err := LoadInstallHosts(optConfigDir)
	if err != nil {
		return err
	}
	var ihs2 []InstallHost
	found := false
	for _, ih := range ihs {
		if ih.Name == name {
			if ih, err = newIh(name); err != nil {
				return err
			}
			found = true
		}
		ihs2 = append(ihs2, ih)
	}
	if !found {
		var ih InstallHost
		if ih, err = newIh(name); err != nil {
			return err
		}
		ihs2 = append(ihs2, ih)
	}
	return StoreInstallHosts(optConfigDir, ihs2)
}
