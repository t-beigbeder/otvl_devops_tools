package common

import (
	"encoding/json"
	"os"
)

func JsonStore(path string, val interface{}) error {
	y, err := json.Marshal(val)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(y)
	if err != nil {
		return err
	}
	return nil
}
