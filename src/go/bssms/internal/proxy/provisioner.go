package proxy

import (
	"bssms/bssms"
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

func handlePrCmds(pc *prCtxType, cmd string) error {
	var err error
	if !pc.established {
		if err = handlePrInitCmd(pc.sbr, cmd, &pc.ins); err != nil {
			return err
		}
		getLogger().Debug("handlePrInitCmd", "ins", pc.ins)
		if err = pc.lner.provisionerReadyEvent(pc.cid, pc.ins); err != nil {
			return err
		}
	} else {
		if err = handlePrInstallCmd(pc.sbr, cmd, &pc.in); err != nil {
			return err
		}
		getLogger().Debug("handlePrInstallCmd", "in", pc.in)
		if err = pc.lner.provisionerInstallEvent(pc.cid, pc.in); err != nil {
			return err
		}
	}
	pc.established = true
	return nil
}
