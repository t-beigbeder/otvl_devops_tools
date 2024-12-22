package common

import "os/exec"

func ExecuteSimple(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	return cmd.Run()
}
