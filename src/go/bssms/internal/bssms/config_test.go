//go:build unit

package bssms

import (
	"bssms/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	td := t.TempDir()
	require.Equal(t, common.AsInterfaceSlice(td, nil), common.AsInterfaceSlice(GetConfigDir(td)))
	require.Equal(t, common.AsInterfaceSlice(td+"/cd", nil), common.AsInterfaceSlice(GetConfigDir(td+"/cd")))
	common.Putenv(BssmsPathEnv, td+"/bp")
	require.Equal(t, common.AsInterfaceSlice(td+"/bp", nil), common.AsInterfaceSlice(GetConfigDir("")))
	common.Delenv(BssmsPathEnv)
	require.Equal(t, common.AsInterfaceSlice(common.Getenv("HOME")+"/.config/.bssms", nil), common.AsInterfaceSlice(GetConfigDir("")))
}
