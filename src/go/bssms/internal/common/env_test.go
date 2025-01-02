//go:build unit

package common

import (
	"testing"
)

func TestGetenv(t *testing.T) {
	require.Equal(t, "", Getenv("BSSMS_ENV"))
	Putenv("BSSMS_ENV", "value")
	require.Equal(t, "value", Getenv("BSSMS_ENV"))
	Delenv("BSSMS_ENV")
	require.Equal(t, "", Getenv("BSSMS_ENV"))
}
