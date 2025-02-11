package lik

import "testing"

func TestEncode(t *testing.T) {
	data := `{"str":"str"}`
	item := ItemFromString(data)
	if item == nil {
		t.Errorf("ERROR of ItemFromString()")
	}
}
