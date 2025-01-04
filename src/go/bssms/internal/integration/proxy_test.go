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

func TestRunCollaboration1(t *testing.T) {
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

	inCancel, err := RunTestInstaller(t, dataFile, 0, proxyPort)
	require.NoError(t, err)
	require.NotNil(t, inCancel)
	defer inCancel()

	time.Sleep(200 * time.Millisecond)
	pxCancel()
	time.Sleep(100 * time.Millisecond)

}

func TestRunCollaboration2(t *testing.T) {
	const (
		proxyPort = "11443"
		dataFile  = "ih2.yaml"
	)
	pxCancel, err := RunTestProxy(proxyPort)
	require.NoError(t, err)
	require.NotNil(t, pxCancel)
	defer pxCancel()

	inCancel0, err := RunTestInstaller(t, dataFile, 0, proxyPort)
	require.NoError(t, err)
	require.NotNil(t, inCancel0)
	defer inCancel0()

	prCancel, err := RunTestProvisioner(dataFile, proxyPort)
	require.NoError(t, err)
	require.NotNil(t, prCancel)
	defer prCancel()

	inCancel1, err := RunTestInstaller(t, dataFile, 1, proxyPort)
	require.NoError(t, err)
	require.NotNil(t, inCancel1)
	defer inCancel1()

	time.Sleep(200 * time.Millisecond)
	pxCancel()
	time.Sleep(100 * time.Millisecond)

}
