package node

import (
	"strings"

	"github.com/PUMATeam/catapult/model"
	"github.com/PUMATeam/catapult/util"
	uuid "github.com/satori/go.uuid"
)

// NodeService exposes operations to perform on a host
type NodeService interface {
	InstallHost(host model.Host) error
	ListVMs() ([]uuid.UUID, error)
	StartVM(vmId uuid.UUID) error
	StopVM(vmId uuid.UUID) error
}

type Node struct {
	Host model.Host
}

// NewNode creates a Node instance
func NewNode(host model.Host) *Node {
	return &Node{
		Host: host,
	}
}

// InstallHost installs prerequisits on the host
func (n *Node) InstallHost() error {
	hi := hostInstall{
		User:            n.Host.User,
		FcVersion:       FcVersion,
		AnsiblePassword: n.Host.Password,
	}
	ac := util.NewAnsibleCommand(util.SetupHostPlaybook,
		hi.User,
		n.Host.Address,
		util.StructToMap(hi, strings.ToLower))
	err := ac.ExecuteAnsible()
	if err != nil {
		return err
	}

	return nil
}

// TODO make it configurable
const FcVersion = "0.15.0"

type hostInstall struct {
	User            string `json:"ignore"`
	AnsiblePassword string `json:"ansible_ssh_pass"`
	FcVersion       string
}
