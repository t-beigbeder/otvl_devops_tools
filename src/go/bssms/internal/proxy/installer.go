package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bufio"
	"fmt"
)

func handleInCmd(sbr *bufio.Reader, cmd string, in *bssms.Installable) error {
	if cmd != bssms.InstallerInstallable {
		return fmt.Errorf("unknown command %s", cmd)
	}
	if err := common.ReadJsonFromStream(sbr, in); err != nil {
		return err
	}
	return nil
}
