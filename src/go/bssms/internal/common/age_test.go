package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAgeEncDec(t *testing.T) {
	pub, pri, err := NewKeyPair()
	require.NoError(t, err)
	ebs, err := EncryptMsg("TestAgeEncDec", pub)
	require.NoError(t, err)
	dbs, err := DecryptMsg(ebs, pri)
	require.NoError(t, err)
	require.Equal(t, "TestAgeEncDec", dbs)
}
