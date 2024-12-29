package proxy

import (
	"bssms/internal/bssms"
	"bssms/internal/common"
	"bufio"
	"fmt"
	"github.com/quic-go/quic-go"
)

func handlePrCmd(stream quic.Stream, cmd string) ([]bssms.Installable, error) {
	sr := bufio.NewReaderSize(stream, common.CtrlDataMaxLn)
	if cmd != bssms.ProvisionerInstallables {
		return nil, fmt.Errorf("unknown command %s", cmd)
	}
	ins := []bssms.Installable{}
	if err := common.ReadJsonFromStream(sr, &ins); err != nil {
		return nil, err
	}
	getLogger().Debug("handlePrCmd", "ins", ins)
	return ins, nil
}
