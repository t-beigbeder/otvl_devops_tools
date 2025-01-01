package proxy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeListener(t *testing.T) {
	lner, err := makeListener()
	assert.NoError(t, err)
	lner.close()
}
