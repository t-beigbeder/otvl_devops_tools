package provisioner

import (
	"bssms/internal/common"
	"testing"
)

func TestGetOSServers(t *testing.T) {
	common.SkipUnlessEnv(t, "BSSMS_TEST_OS")
	ihs, err := GetOSServers()
	if err != nil {
		t.Fatal(err)
	}
	_ = ihs
}
