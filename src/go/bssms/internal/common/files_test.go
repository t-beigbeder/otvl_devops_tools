package common

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDir(t *testing.T) {
	tp := t.TempDir()
	require.Equal(t, AsInterfaceSlice(false, nil), AsInterfaceSlice(CheckDir(filepath.Join(tp, "nocreate"), false, false, 0o777)))
	require.Equal(t, AsInterfaceSlice(true, nil), AsInterfaceSlice(CheckDir(filepath.Join(tp, "create"), true, false, 0o777)))
	r, err := CheckDir(filepath.Join(tp, "level1/nocreateall"), true, false, 0o777)
	require.Equal(t, false, r)
	require.Error(t, err)
	require.Equal(t, AsInterfaceSlice(true, nil), AsInterfaceSlice(CheckDir(filepath.Join(tp, "level1/createall"), true, true, 0o777)))
	l1l2 := filepath.Join(tp, "level1/level2")
	require.Equal(t, AsInterfaceSlice(true, nil), AsInterfaceSlice(CheckDir(l1l2, true, true, 0o500)))
	if os.Getuid() == 0 {
		require.NoError(t, ExecuteSimple("chattr", "+i", l1l2))
	}
	r, err = CheckDir(filepath.Join(tp, "level1/level2/refused"), true, false, 0o777)
	require.Equal(t, false, r)
	require.Error(t, err)
	if os.Getuid() == 0 {
		require.NoError(t, ExecuteSimple("chattr", "-i", l1l2))
	}
}
