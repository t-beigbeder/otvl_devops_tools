package bssms

import "testing"

func TestGetIPPort(t *testing.T) {
	i, p, err := GetIPPort(":443")
	if err != nil || i != nil || p != 443 {
		t.Error(i, p, err)
	}
	i, p, err = GetIPPort("localhost:443")
	if err != nil || i != nil || p != 443 {
		t.Error(i, p, err)
	}
	i, p, err = GetIPPort("0.0.0.0:443")
	if err != nil || i == nil || p != 443 {
		t.Error(i, p, err)
	}
}
