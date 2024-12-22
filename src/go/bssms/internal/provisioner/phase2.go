package provisioner

func RunPhase2(optConfigDir string) error {
	var ihs0 []InstallHost
	var err error
	if ihs0, err = LoadInstallHosts(optConfigDir); err != nil {
		return err
	}
	for _, ihs := range ihs0 {
		_ = ihs
	}
	return nil
}
