package bssms

import "bssms/internal/bssms"

type InstallHost struct {
	bssms.Installable
	PrivateKey string
	PrPubKey   string
}
