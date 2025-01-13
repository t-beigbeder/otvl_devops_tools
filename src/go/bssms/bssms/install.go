package bssms

type Installable struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	ServerUuid    string `json:"serverUuid,omitempty" yaml:"serverUuid,omitempty"`
	MacExtAddress string `json:"macExtAddress,omitempty" yaml:"macExtAddress,omitempty"`
	MacIntAddress string `json:"macIntAddress,omitempty" yaml:"macIntAddress,omitempty"`
	MacAddress    string `json:"macAddress,omitempty" yaml:"macAddress,omitempty"`
	IPExtAddress  string `json:"IPExtAddress,omitempty" yaml:"IPExtAddress,omitempty"`
	IPIntAddress  string `json:"IPIntAddress,omitempty" yaml:"IPIntAddress,omitempty"`
	IPAddress     string `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
	EncSecrets    string `json:"encSecrets,omitempty" yaml:"encSecrets,omitempty"`
	Installing    bool   `json:"-" yaml:"-"`
	Installed     bool   `json:"-" yaml:"-"`
}

func (iin Installable) Matches(pin Installable) bool {
	if pin.ServerUuid == iin.ServerUuid &&
		(pin.MacIntAddress == iin.MacAddress || pin.MacExtAddress == iin.MacAddress || pin.MacAddress == iin.MacAddress) &&
		(pin.IPIntAddress == iin.IPAddress || pin.IPExtAddress == iin.IPAddress || pin.IPAddress == iin.IPAddress) {
		return true
	}
	return false
}

type Secrets map[string]string
