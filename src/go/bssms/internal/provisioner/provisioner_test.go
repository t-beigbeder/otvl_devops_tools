package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/integration"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"path/filepath"
	"testing"
	"time"
)

func getIhs(dataFile string) ([]InstallHost, error) {
	return LoadYamlInstallHosts(filepath.Join(integration.GetTestDataDir(), dataFile))
}

func TestProvisionerRun(t *testing.T) {
	pxCancel, err := integration.RunProxy(t)
	assert.NoError(t, err)
	defer pxCancel()
	prCtx, prCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer prCancel()
	ihs, err := getIhs("ih1.yaml")
	assert.NoError(t, err)
	err = run(&bssms.ProvisionerConfig{
		BaseConfig:   bssms.BaseConfig{Ctx: prCtx},
		UnsafeTls:    true,
		ProxyAddress: integration.ProxyAddress(),
	}, ihs)
	assert.NoError(t, err)
}
