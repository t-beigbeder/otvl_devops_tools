package provisioner

import "bssms/internal/bssms"

func RunPhase0(ss []string) error {
	var ihs []InstallHost
	for _, s := range ss {
		ihs = append(ihs, InstallHost{Installable: bssms.Installable{Name: s}})
	}

}
