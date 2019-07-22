package node

import (
	"github.com/PUMATeam/catapult/model"
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
	host model.Host
}

// InstallHost installs prerequisits on the host
func (n *Node) InstallHost() {

}

type HostInstall struct {
	Address  string
	User     string
	Password string
}
