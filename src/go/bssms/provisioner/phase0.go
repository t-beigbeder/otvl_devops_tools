package provisioner

import (
	"bssms/bssms"
	"bssms/internal/common"
)

func RunPhase0(optConfigDir string, ss []string) error {
	var ihs []InstallHost
	for _, s := range ss {
		pub, pri, err := common.NewKeyPair()
		if err != nil {
			return err
		}
		ihs = append(ihs, InstallHost{
			Installable: bssms.Installable{Name: s},
			PubKey:      pub,
			PrivateKey:  pri,
		})
	}
	return StoreInstallHosts(optConfigDir, ihs)
}
