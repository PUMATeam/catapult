package util

import (
	"fmt"
)

const (
	SetupHostPlaybook  = "ansible/roles/setup_host/playbook.yml"
	ansiblePlaybookCmd = "ansible-playbook"
)

type AnsibleCommand struct {
	Playbook    string
	ExtraParams map[string]string
	User        string
	Host        string
	cmd         []string
}

// NewAnsibleCommand creates a new AnsibleCommand
func NewAnsibleCommand(playbook, user, host string, params map[string]string) *AnsibleCommand {
	ac := &AnsibleCommand{
		Playbook:    playbook,
		ExtraParams: params,
		User:        user,
		Host:        host,
	}

	ac.generateCmd()

	return ac
}

func (ac *AnsibleCommand) generateCmd() {
	cmd := []string{
		ac.Playbook,
		"-i",
		fmt.Sprintf("%s,", ac.Host),
		"-u",
		ac.User,
	}

	for key, value := range ac.ExtraParams {
		cmd = append(cmd, fmt.Sprintf("-e %s=%s", key, value))
	}

	ac.cmd = cmd
}

// ExecuteAnsible executes a given ansible playbook
func (ac *AnsibleCommand) ExecuteAnsible() error {
	return ExecuteCmd(ansiblePlaybookCmd, ac.cmd)
}
