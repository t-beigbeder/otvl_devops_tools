package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bufio"
	"fmt"
)

func handlePrInitCmd(sbr *bufio.Reader, cmd string, ins *[]bssms.Installable) error {
	if cmd != bssms.ProvisionerInstallables {
		return fmt.Errorf("unknown command %s", cmd)
	}
	if *ins != nil {
		return fmt.Errorf("provisioner already sent %d installables", len(*ins))
	}
	if err := common.ReadJsonFromStream(sbr, ins); err != nil {
		return err
	}
	return nil
}

func handlePrInstallCmd(sbr *bufio.Reader, cmd string, in *bssms.Installable) error {
	if cmd != bssms.ProvisionerInstall {
		return fmt.Errorf("unknown command %s", cmd)
	}
	if err := common.ReadJsonFromStream(sbr, in); err != nil {
		return err
	}
	return nil
}
