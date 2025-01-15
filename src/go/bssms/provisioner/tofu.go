package provisioner

import (
	"bssms/bssms"
)

func TofuRun(config *bssms.ProvisionerConfig, optConfigDir string, ih InstallHost) error {
	_, err := ReadInstallHost(optConfigDir, ih.Name)
	if err != nil {
		return err
	}
	err = RunIhs(config, []*InstallHost{&ih})
	if err != nil {
		return err
	}
	err = SaveInstallHost(optConfigDir, ih)
	if err != nil {
		return err
	}
	return nil
}
