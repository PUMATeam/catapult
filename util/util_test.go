package util

import (
	"reflect"
	"strings"
	"testing"

	node "github.com/PUMATeam/catapult/node"
)

func TestStructToMapLowerCase(t *testing.T) {
	h := node.HostInstall{
		Address:  "192.168.1.1",
		User:     "vagrant",
		Password: "vagrant",
	}

	got := StructToMap(h, strings.ToLower)
	expected := map[string]string{
		"address":  "192.168.1.1",
		"user":     "vagrant",
		"password": "vagrant",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\n\tGOT: %s \n\tEXPECTED: %s", got, expected)
	}
}
