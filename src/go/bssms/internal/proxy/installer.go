package proxy

import (
	"bssms/bssms"
	"bssms/internal/common"
	"bufio"
	"fmt"
)

func handleInInstallableCmd(sbr *bufio.Reader, cmd string, in *bssms.Installable) error {
	if cmd != bssms.InstallerInstallable {
		return fmt.Errorf("unknown command %s", cmd)
	}
	if err := common.ReadJsonFromStream(sbr, in); err != nil {
		return err
	}
	return nil
}

func handleInInstalledCmd(sbr *bufio.Reader, cmd string, in *bssms.Installable) error {
	if cmd != bssms.InstallerInstalled {
		return fmt.Errorf("unknown command %s", cmd)
	}
	if err := common.ReadJsonFromStream(sbr, in); err != nil {
		return err
	}
	return nil
}

func handleInCmds(ic *inCtxType, cmd string) error {
	var err error
	if !ic.established {
		if err = handleInInstallableCmd(ic.sbr, cmd, &ic.in); err != nil {
			return err
		}
		getLogger().Debug("handleInInstallableCmd", "in", ic.in)
		if err = ic.lner.installerReadyEvent(ic.cid, ic.in); err != nil {
			return err
		}
	} else {
		if err = handleInInstalledCmd(ic.sbr, cmd, &ic.in); err != nil {
			return err
		}
		getLogger().Debug("handleInInstalledCmd", "in", ic.in)
		if err = ic.lner.installerInstalledEvent(ic.cid, ic.in); err != nil {
			return err
		}
	}
	ic.established = true
	return nil
}
