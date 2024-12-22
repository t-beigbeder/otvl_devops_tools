package provisioner

import (
	"os"
	"testing"
)

func TestGetOSServers(t *testing.T) {
	if os.Getenv("BSSMS_TEST_OS") == "" {
		t.Skip("BSSMS_TEST_OS not set")
	}
	ihs, err := GetOSServers()
	if err != nil {
		t.Fatal(err)
	}
	_ = ihs
}
