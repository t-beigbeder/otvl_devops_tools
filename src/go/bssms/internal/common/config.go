package common

import (
	"path/filepath"
)

func SystemConfigDir() string {
	var scd string
	if scd = Getenv("XDG_CONFIG_HOME"); scd == "" {
		return filepath.Join(Getenv("HOME"), ".config")
	}
	return scd
}
