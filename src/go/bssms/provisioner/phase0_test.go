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

func TestMergePhase0(t *testing.T) {
	td := t.TempDir()
	if err := RunPhase0(td, []string{"tbst6", "h1", "h2", "h3"}); err != nil {
		t.Fatal(err)
	}
	ihs, err := LoadInstallHosts(td)
	require.NoError(t, err)
	require.Equal(t, len(ihs), 4)
	require.NotNil(t, ihs[0].PrivateKey)
	err = MergePhase0(td, "h4")
	require.NoError(t, err)
	ihs, err = LoadInstallHosts(td)
	require.NoError(t, err)
	require.Equal(t, len(ihs), 5)
	err = MergePhase0(td, "h2")
	require.NoError(t, err)
	ihs2, err := LoadInstallHosts(td)
	require.NoError(t, err)
	require.Equal(t, len(ihs2), 5)
	require.Equal(t, ihs[1].PrivateKey, ihs2[1].PrivateKey)
	require.NotEqual(t, ihs[2].PrivateKey, ihs2[2].PrivateKey)
}
