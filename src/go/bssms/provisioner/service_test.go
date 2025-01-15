package provisioner

import (
	"bssms/bssms"
	"github.com/stretchr/testify/require"
	"testing"
)

func testIhs() []InstallHost {
	return []InstallHost{
		{
			Installable: bssms.Installable{
				Name: "ih1",
			},
			PrivateKey: "prik1",
			PubKey:     "pubk1",
			Secrets:    map[string]string{"key1": "val1", "key2": "val2"},
		},
		{
			Installable: bssms.Installable{
				Name: "ih2",
			},
			PrivateKey: "prik2",
			PubKey:     "pubk2",
		},
	}
}

func TestStoreAndLoadInstallHosts(t *testing.T) {
	tihs := testIhs()
	td := t.TempDir()
	err := StoreInstallHosts(td, tihs)
	require.NoError(t, err)
	ihs2, err := LoadInstallHosts(td)
	require.NoError(t, err)
	require.Equal(t, tihs, ihs2)
	ihs3, err := LoadFilteredInstallHosts(td, nil)
	require.NoError(t, err)
	require.Equal(t, tihs, ihs3)
	ihs4, err := LoadFilteredInstallHosts(td, []string{"ih2", "ih1"})
	require.NoError(t, err)
	require.Equal(t, tihs, ihs4)
	ihs5, err := LoadFilteredInstallHosts(td, []string{"ih2"})
	require.NoError(t, err)
	require.Equal(t, tihs[1:2], ihs5)
}

func TestReadAndSaveInstallHost(t *testing.T) {
	tihs := testIhs()
	td := t.TempDir()
	err := StoreInstallHosts(td, tihs)
	require.NoError(t, err)
	tih, err := ReadInstallHost(td, "ih3")
	require.NoError(t, err)
	require.Equal(t, tih.Name, "")
	ih3 := InstallHost{
		Installable: bssms.Installable{
			Name: "ih3",
		},
		PrivateKey: "prik3",
		PubKey:     "pubk3",
	}
	err = SaveInstallHost(td, ih3)
	require.NoError(t, err)
	ihs, err := LoadInstallHosts(td)
	require.NoError(t, err)
	require.Equal(t, 3, len(ihs))
	tih, err = ReadInstallHost(td, "ih3")
	require.NoError(t, err)
	require.Equal(t, tih, ih3)
	ih3.Secrets = map[string]string{"key3a": "val3a"}
	err = SaveInstallHost(td, ih3)
	require.NoError(t, err)
	require.Equal(t, 3, len(ihs))
	tih, err = ReadInstallHost(td, "ih3")
	require.NoError(t, err)
	require.Equal(t, tih, ih3)
}
