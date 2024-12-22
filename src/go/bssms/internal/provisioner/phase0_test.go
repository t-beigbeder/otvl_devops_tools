package provisioner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunPhase0(t *testing.T) {
	td := t.TempDir()
	if err := RunPhase0(td, []string{"h1", "h2", "h3"}); err != nil {
		t.Fatal(err)
	}
	ihs, err := LoadInstallHosts(td)
	assert.NoError(t, err)
	assert.Equal(t, len(ihs), 3)
	assert.NotNil(t, ihs[0].PrivateKey)
}
