package util

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	SetupHostPlaybook    = "ansible/roles/setup_host/playbook.yml"
	ActivateHostPlaybook = "ansible/roles/setup_host/activate_playbook.yml"
	ansiblePlaybookCmd   = "ansible-playbook"
)

type AnsibleCommand struct {
	Playbook    string
	ExtraParams map[string]string
	User        string
	Host        string
	cmd         []string
	log         *log.Logger
}

// NewAnsibleCommand creates a new AnsibleCommand
func NewAnsibleCommand(playbook, user, host string, params map[string]string, logger *log.Logger) *AnsibleCommand {
	ac := &AnsibleCommand{
		Playbook:    playbook,
		ExtraParams: params,
		User:        user,
		Host:        host,
		log:         logger,
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

	ac.log.WithFields(log.Fields{
		"command": cmd,
	}).Info("Generated ansible command")

	ac.cmd = cmd
}

// ExecuteAnsible executes a given ansible playbook
func (ac *AnsibleCommand) ExecuteAnsible() error {
	return ExecuteCmd(ansiblePlaybookCmd, ac.cmd)
}
