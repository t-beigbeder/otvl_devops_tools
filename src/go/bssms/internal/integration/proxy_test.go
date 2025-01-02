package integration

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRunProxy(t *testing.T) {
	const proxyPort = "9443"
	cancel, err := RunTestProxy(proxyPort)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestRunCollaboration(t *testing.T) {
	const (
		proxyPort = "10443"
		dataFile  = "ih1.yaml"
	)
	pxCancel, err := RunTestProxy(proxyPort)
	require.NoError(t, err)
	require.NotNil(t, pxCancel)
	defer pxCancel()

	prCancel, err := RunTestProvisioner(dataFile, proxyPort)
	require.NoError(t, err)
	require.NotNil(t, prCancel)
	defer prCancel()

	inCancel, err := RunTestInstaller(dataFile, proxyPort)
	require.NoError(t, err)
	require.NotNil(t, inCancel)
	defer inCancel()

	time.Sleep(100 * time.Millisecond)
	pxCancel()
	time.Sleep(100 * time.Millisecond)

}
