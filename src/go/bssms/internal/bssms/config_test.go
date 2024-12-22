//go:build unit

package bssms

import (
	"bssms/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	td := t.TempDir()
	assert.Equal(t, common.AsInterfaceSlice(td, nil), common.AsInterfaceSlice(GetConfigDir(td)))
	assert.Equal(t, common.AsInterfaceSlice(td+"/cd", nil), common.AsInterfaceSlice(GetConfigDir(td+"/cd")))
	common.Putenv(BssmsPathEnv, td+"/bp")
	assert.Equal(t, common.AsInterfaceSlice(td+"/bp", nil), common.AsInterfaceSlice(GetConfigDir("")))
	common.Delenv(BssmsPathEnv)
	assert.Equal(t, common.AsInterfaceSlice(common.Getenv("HOME")+"/.config/.bssms", nil), common.AsInterfaceSlice(GetConfigDir("")))
}
