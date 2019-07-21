package util

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

const (
	SetupHostPlaybook  = "ansible/roles/setup_host/playbook.yml"
	ansiblePlaybookCmd = "ansible-playbook"
)

func ExecuteAnsible(playbook string, host string) {
	hostParam := fmt.Sprintf("%s,", host)
	userParam := "vagrant"
	sshKeyParam := "--private-key=./tests/insecure_private_key"
	fcVersionParam := "fc_version=0.15.0"
	ansibleParams := fmt.Sprintf("-e \"user=%s %s %s\"", userParam, sshKeyParam, fcVersionParam)
	var outb, errb bytes.Buffer
	cmd := exec.Command(ansiblePlaybookCmd,
		playbook,
		"-i",
		hostParam,
		"-u",
		userParam,
		sshKeyParam,
		ansibleParams,
		"-v",
	)

	log.Printf("Running %s with %s", cmd.Path, cmd.Args)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		log.Println(errb.String())
	}

	fmt.Println("out:", outb.String())
}
