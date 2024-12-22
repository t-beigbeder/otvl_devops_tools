package provisioner

import (
	"bssms/internal/bssms"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, tihs, ihs2)
}
