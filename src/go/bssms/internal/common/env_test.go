//go:build unit

package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetenv(t *testing.T) {
	assert.Equal(t, "", Getenv("BSSMS_ENV"))
	Putenv("BSSMS_ENV", "value")
	assert.Equal(t, "value", Getenv("BSSMS_ENV"))
	Delenv("BSSMS_ENV")
	assert.Equal(t, "", Getenv("BSSMS_ENV"))
}
