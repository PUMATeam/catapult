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
	userParam := "user=vagrant"
	ansibleUserParam := "ansible_user=vagrant"
	sshKeyParam := "ansible_ssh_private_key_file=./tests/insecure_private_key"
	fcVersionParam := "fc_version=0.15.0"
	ansibleParams := fmt.Sprintf("-e \"%s %s %s %s\"", userParam, ansibleUserParam, sshKeyParam, fcVersionParam)
	var outb, errb bytes.Buffer
	cmd := exec.Command(ansiblePlaybookCmd,
		playbook,
		"-i",
		hostParam,
		ansibleParams,
	)
	log.Printf("Running %s with %s", cmd.Path, cmd.Args)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	fmt.Println("out:", outb.String(), "err:", errb.String())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("out:", outb.String(), "err:", errb.String())
}
