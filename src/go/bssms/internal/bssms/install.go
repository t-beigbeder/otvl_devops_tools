package bssms

type Installable struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	ServerUuid string `json:"serverUuid,omitempty" yaml:"serverUuid,omitempty"`
	MacAddress string `json:"macAddress,omitempty" yaml:"macAddress,omitempty"`
	IPAddress  string `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
}
