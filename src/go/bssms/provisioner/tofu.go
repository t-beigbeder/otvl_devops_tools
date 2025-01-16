package provisioner

import (
	"bssms/bssms"
)

func TofuRun(config *bssms.ProvisionerConfig, optConfigDir string, ih InstallHost) error {
	_, err := ReadInstallHost(optConfigDir, ih.Name)
	if err != nil {
		return err
	}
	ih.TofuRunning = len(ih.Secrets) > 0
	ih.TofuError = ""
	ih.Installed = len(ih.Secrets) == 0
	err = SaveInstallHost(optConfigDir, ih)
	if err != nil {
		return err
	}
	if len(ih.Secrets) == 0 {
		return nil
	}

	err = RunIhs(config, []*InstallHost{&ih})
	ih.TofuRunning = false
	if err != nil {
		ih.TofuError = err.Error()
		_ = SaveInstallHost(optConfigDir, ih)
		return err
	}
	err = SaveInstallHost(optConfigDir, ih)
	if err != nil {
		return err
	}
	return nil
}
