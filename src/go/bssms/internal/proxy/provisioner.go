package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bufio"
	"fmt"
)

func handlePrCmd(streamReader *bufio.Reader, cmd string) ([]bssms.Installable, error) {
	if cmd != bssms.ProvisionerInstallables {
		return nil, fmt.Errorf("unknown command %s", cmd)
	}
	ins := []bssms.Installable{}
	if err := common.ReadJsonFromStream(streamReader, &ins); err != nil {
		return nil, err
	}
	getLogger().Debug("handlePrCmd", "ins", ins)
	return ins, nil
}
