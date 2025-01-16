package provisioner

import (
	"bssms/bssms"
	"bssms/internal/common"
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

func MergePhase0(optConfigDir string, name string) error {
	ih, err := newIh(name)
	if err != nil {
		return err
	}
	return SaveInstallHost(optConfigDir, ih)
}
