package provisioner

import (
	"bssms/internal/bssms"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreAndLoadInstallHosts(t *testing.T) {
	tihs := []InstallHost{
		{
			Installable: bssms.Installable{
				Name: "ih1",
			},
			PrivateKey: "prik1",
			PubKey:     "pubk1",
		},
		{
			Installable: bssms.Installable{
				Name: "ih2",
			},
			PrivateKey: "prik2",
			PubKey:     "pubk2",
		},
	}
	td := t.TempDir()
	err := StoreInstallHosts(td, tihs)
	if err != nil {
		t.Fatal(err)
	}
	ihs2, err := LoadInstallHosts(td)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tihs, ihs2)
	ihs3, err := LoadFilteredInstallHosts(td, nil)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tihs, ihs3)
	ihs4, err := LoadFilteredInstallHosts(td, []string{"ih2", "ih1"})
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tihs, ihs4)
	ihs5, err := LoadFilteredInstallHosts(td, []string{"ih2"})
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tihs[1:2], ihs5)
}
