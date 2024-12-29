package bssms

type Installable struct {
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	ServerUuid   string `json:"serverUuid,omitempty" yaml:"serverUuid,omitempty"`
	MacAddress   string `json:"macAddress,omitempty" yaml:"macAddress,omitempty"`
	IPExtAddress string `json:"IPExtAddress,omitempty" yaml:"IPExtAddress,omitempty"`
	IPIntAddress string `json:"IPIntAddress,omitempty" yaml:"IPIntAddress,omitempty"`
}
