package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewAnsibleCommand(t *testing.T) {
	ac := NewAnsibleCommand(SetupHostPlaybook,
		"root",
		"192.168.50.50",
		map[string]string{
			"ansible_ssh_pass": "centos",
			"fc_version":       "0.15.0",
		},
	)

	got := strings.Join(ac.cmd, " ")
	expected := "ansible/roles/setup_host/playbook.yml -i 192.168.50.50, -u root -e ansible_ssh_pass=centos -e fc_version=0.15.0"
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("\n\tGOT: %s \n\tEXPECTED: %s", got, expected)
	}
}
