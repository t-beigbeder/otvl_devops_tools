package provisioner

import (
	"bssms/internal/common"
	"testing"
)

func TestRunPhase2(t *testing.T) {
	common.SkipUnlessEnv(t, "BSSMS_TEST_OS")
	td := t.TempDir()
	if err := RunPhase0(td, []string{"tbst6", "hst2"}); err != nil {
		t.Fatal(err)
	}
	err := RunPhase2(td, []string{})
	if err != nil {
		return
	}
}
