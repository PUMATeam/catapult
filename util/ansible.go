package util

const (
	SetupHostPlaybook  = "ansible/roles/setup_host/playbook.yml"
	ansiblePlaybookCmd = "ansible-playbook"
)

type AnsibleCommand struct {
	Playbook    string
	ExtraParams map[string]string
	User        string
}

func generateCmd(playbook string, extraParams map[string]string) string {
	return ""
}

// func ExecuteAnsible(params map[string]string) {
// 	hostParam := fmt.Sprintf("%s,", host)
// 	userParam := "vagrant"
// 	hostUserParam := fmt.Sprintf("host_user=%s", userParam)
// 	sshKeyParam := "--private-key=./tests/insecure_private_key"
// 	fcVersionParam := "fc_version=0.15.0"
// 	extraVars := fmt.Sprintf("%s %s", hostUserParam, fcVersionParam)
// 	ansibleParams := fmt.Sprintf("-e %s", extraVars)
// 	var outb, errb bytes.Buffer
// 	cmd := exec.Command(ansiblePlaybookCmd,
// 		playbook,
// 		"-i",
// 		hostParam,
// 		"-u",
// 		userParam,
// 		sshKeyParam,
// 		ansibleParams,
// 		"-v",
// 	)

// 	log.Printf("Running with %s", cmd.Args)
// 	cmd.Stdout = &outb
// 	cmd.Stderr = &errb

// 	err := cmd.Run()
// 	if err != nil {
// 		log.Println(errb.String())
// 	}

// 	fmt.Println("out:", outb.String())
// }
