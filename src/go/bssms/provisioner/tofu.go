package provisioner

import (
	"bssms/bssms"
)

func TofuRun(config *bssms.ProvisionerConfig, optConfigDir string, ih InstallHost) error {
	_, err := ReadInstallHost(optConfigDir, ih.Name)
	if err != nil {
		return err
	}
	ih.TofuRunning = true
	err = SaveInstallHost(optConfigDir, ih)
	if err != nil {
		return err
	}
	err = RunIhs(config, []*InstallHost{&ih})
	ih.TofuRunning = false
	if err != nil {
		_ = SaveInstallHost(optConfigDir, ih)
		return err
	}
	err = SaveInstallHost(optConfigDir, ih)
	if err != nil {
		return err
	}
	return nil
}
