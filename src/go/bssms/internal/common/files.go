package common

import "os"

func CheckDir(path string, create, all bool, perm os.FileMode) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return true, nil
	}
	if err != nil && os.IsNotExist(err) {
		if !create {
			return false, nil
		}
		if !all {
			err = os.Mkdir(path, perm)
		} else {
			err = os.MkdirAll(path, perm)
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return true, nil
}
