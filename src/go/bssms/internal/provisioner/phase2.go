package provisioner

import (
	"bssms/internal/bssms"
)

func RunPhase2(optConfigDir string) error {
	var ihs0, ihs1, ihs2 []InstallHost
	var err error
	if ihs0, err = LoadInstallHosts(optConfigDir); err != nil {
		return err
	}
	if ihs1, err = GetOSServers(); err != nil {
		return err
	}
	for _, ih0 := range ihs0 {
		for _, ih1 := range ihs1 {
			if ih1.Name == ih0.Name {
				ih2 := InstallHost{
					Installable: bssms.Installable{
						Name:         ih0.Name,
						ServerUuid:   ih1.ServerUuid,
						MacAddress:   ih1.MacAddress,
						IPExtAddress: ih1.IPExtAddress,
						IPIntAddress: ih1.IPIntAddress,
					},
					PrivateKey: ih0.PrivateKey,
					PubKey:     ih0.PubKey,
				}
				ihs2 = append(ihs2, ih2)
				break
			}
		}
	}
	return StoreInstallHosts(optConfigDir, ihs2)
}
