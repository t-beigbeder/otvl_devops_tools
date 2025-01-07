package installer

import (
	"bssms/bssms"
)

type InstallHost struct {
	bssms.Installable
	PrivateKey string
	PrPubKey   string
}
