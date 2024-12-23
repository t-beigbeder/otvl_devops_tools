package integration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunProxy(t *testing.T) {
	cancel, err := RunProxy(t)
	assert.NoError(t, err)
	cancel()
}
