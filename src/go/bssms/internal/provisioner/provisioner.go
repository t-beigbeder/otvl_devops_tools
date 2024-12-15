package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/qutils"
	"fmt"
	"os"
)

func RunProvisioner(config *bssms.ProvisionerConfig) error {
	// FIXME: 2024/12/15 17:43:20 CRYPTO_ERROR 0x170 (remote): tls: unrecognized name
	// https://github.com/alta/insecure
	conn, err := qutils.GetQuicConn(config.ProxyAddress)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "RunProvisioner %v", conn)
	return nil
}
