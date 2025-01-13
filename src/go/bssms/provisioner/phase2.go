package provisioner

import (
	"bssms/bssms"
)

func RunPhase2(optConfigDir string, ss []string, secf map[string]map[string]string) error {
	var ihs0, ihs1, ihs2 []InstallHost
	var err error
	if ihs0, err = LoadFilteredInstallHosts(optConfigDir, ss); err != nil {
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
						Name:          ih0.Name,
						ServerUuid:    ih1.ServerUuid,
						MacExtAddress: ih1.MacExtAddress,
						MacIntAddress: ih1.MacIntAddress,
						IPExtAddress:  ih1.IPExtAddress,
						IPIntAddress:  ih1.IPIntAddress,
					},
					PrivateKey: ih0.PrivateKey,
					PubKey:     ih0.PubKey,
				}
				if secs, ok := secf[ih0.Name]; ok {
					ih2.Secrets = secs
				}
				ihs2 = append(ihs2, ih2)
				break
			}
		}
	}
	return StoreInstallHosts(optConfigDir, ihs2)
}
