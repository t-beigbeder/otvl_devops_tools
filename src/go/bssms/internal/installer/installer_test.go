package installer

import (
	"bssms/internal/bssms"
	"bssms/internal/integration"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	pxCancel, err := integration.RunProxy(t)
	assert.NoError(t, err)
	defer pxCancel()
	inCtx, prCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer prCancel()
	err = Run(&bssms.InstallerConfig{
		BaseConfig:   bssms.BaseConfig{Ctx: inCtx},
		UnsafeTls:    true,
		ProxyAddress: integration.ProxyAddress(),
		Installable: bssms.Installable{
			ServerUuid: "uuid:de:ad:be:ef",
			MacAddress: "mac:de:ad:be:ef",
			IPAddress:  "127.0.0.42",
		},
	})
	assert.NoError(t, err)
}
