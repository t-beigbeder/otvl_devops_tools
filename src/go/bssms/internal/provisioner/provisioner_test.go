package provisioner

import (
	"bssms/internal/bssms"
	"bssms/internal/integration"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestProvisionerRun(t *testing.T) {
	pxCancel, err := integration.RunProxy(t)
	assert.NoError(t, err)
	defer pxCancel()
	prCtx, prCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer prCancel()
	err = Run(&bssms.ProvisionerConfig{
		BaseConfig:   bssms.BaseConfig{Ctx: prCtx},
		UnsafeTls:    true,
		ProxyAddress: integration.ProxyAddress(),
	})
	assert.NoError(t, err)
}
