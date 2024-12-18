package bssms

type Installable struct {
	ServerUuid string
	MacAddress string
	IPAddress  string
}

type InstallHost struct {
	Installable
	PrivateKey string
	PrPubKey   string
}
