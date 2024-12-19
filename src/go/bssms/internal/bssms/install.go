package bssms

type Installable struct {
	ServerUuid string `json:"serverUuid,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
	IPAddress  string `json:"IPAddress,omitempty"`
}
