package proxy

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestMakeListener(t *testing.T) {
	lner, err := makeListener(context.Background())
	require.NoError(t, err)
	lner.close()
}
