package bssms

type Installable struct {
	Name         string `json:"name,omitempty" yaml:"name,omitempty"`
	ServerUuid   string `json:"serverUuid,omitempty" yaml:"serverUuid,omitempty"`
	MacAddress   string `json:"macAddress,omitempty" yaml:"macAddress,omitempty"`
	IPExtAddress string `json:"IPExtAddress,omitempty" yaml:"IPExtAddress,omitempty"`
	IPIntAddress string `json:"IPIntAddress,omitempty" yaml:"IPIntAddress,omitempty"`
	IPAddress    string `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
}

func (iin Installable) Matches(pin Installable) bool {
	if pin.ServerUuid == iin.ServerUuid &&
		pin.MacAddress == iin.MacAddress &&
		(pin.IPIntAddress == iin.IPAddress || pin.IPExtAddress == iin.IPAddress) {
		return true
	}
	return false
}

type Secrets map[string]string
