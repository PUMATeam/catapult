package util

import (
	"reflect"
	"testing"
)

func TestNewAnsibleCommand(t *testing.T) {
	ac := NewAnsibleCommand(SetupHostPlaybook,
		"vagrant",
		"192.168.50.50",
		map[string]string{
			"fc_version": "0.15.0",
			"host_user":  "vagrant",
		},
	)

	got := ac.cmd
	expected := "ansible-playbook ansible/roles/setup_host/playbook.yml -i 192.168.50.50, -u vagrant -e \"fc_version=0.15.0 host_user=vagrant\""
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\n\tGOT: %s \n\tEXPECTED: %s", got, expected)
	}
}
