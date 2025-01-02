package provisioner

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunPhase0(t *testing.T) {
	td := t.TempDir()
	if err := RunPhase0(td, []string{"tbst6", "h1", "h2", "h3"}); err != nil {
		t.Fatal(err)
	}
	ihs, err := LoadInstallHosts(td)
	require.NoError(t, err)
	require.Equal(t, len(ihs), 4)
	require.NotNil(t, ihs[0].PrivateKey)
}
