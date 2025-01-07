package bssms

import (
	"bssms/internal/common"
	"fmt"
	"path/filepath"
)

const BssmsPathEnv = "BSSMS_PATH"

func GetConfigDir(optConfigDir string) (string, error) {
	var (
		cd     string
		create bool
		all    bool
		ok     bool
		err    error
	)
	if optConfigDir != "" {
		cd = optConfigDir
		create = true
	} else if bp := common.Getenv(BssmsPathEnv); bp != "" {
		cd = bp
		create = true
	} else {
		cd = filepath.Join(common.SystemConfigDir(), ".bssms")
		create = true
	}
	ok, err = common.CheckDir(cd, create, all, 0o700)
	if err != nil {
		return cd, err
	}
	if !ok {
		return cd, fmt.Errorf("GetConfigDir: cannot create directory %s", cd)
	}
	return cd, nil
}
