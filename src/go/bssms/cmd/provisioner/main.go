package main

import "bssms/internal/provisioner"

func main() {
	err := provisioner.NewP()
	if err != nil {
		return
	}
}
