package provisioner

import "bssms/internal/qutils"

func NewP() error {
	err := qutils.RunServer()
	return err
}
