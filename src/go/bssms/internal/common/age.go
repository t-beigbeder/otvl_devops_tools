package common

import "filippo.io/age"

func NewKeyPair() (string, string, error) {
	xi, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", err
	}
	return xi.Recipient().String(), xi.String(), nil
}
