package node

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/PUMATeam/catapult/pkg/model"
	uuid "github.com/satori/go.uuid"
)

// NodeService exposes operations to perform on a host
type NodeService interface {
	ListVMs() ([]uuid.UUID, error)
	StartVM(ctx context.Context, vm model.VM) error
	StopVM(vmId uuid.UUID) error
}

type Node struct {
	host model.Host
}

// NewNodeService creates a Node instance
func NewNodeService(host model.Host) NodeService {
	return &Node{
		host: host,
	}
}

// StartVM starts an FC VM on the host
func (n *Node) StartVM(ctx context.Context, vm model.VM) error {
	// TODO: make port configurable
	conn, err := grpc.Dial(fmt.Sprintf("%s:8888", n.host.Address), grpc.WithInsecure())
	if err != nil {
		return err
	}

	uuid := &UUID{
		Value: vm.ID.String(),
	}

	vmConfig := &VmConfig{
		VmID:           uuid,
		Memory:         vm.Memory,
		Vcpus:          vm.VCPU,
		KernelImage:    vm.KernelImage,
		RootFileSystem: vm.RootFileSystem,
	}
	client := NewNodeClient(conn)
	resp, err := client.StartVM(ctx, vmConfig)
	if err != nil {
		log.Error("grpc error: ", err)
	}

	log.Debug("grpc response:", resp)

	return nil
}

func (n *Node) ListVMs() ([]uuid.UUID, error) {
	return nil, nil
}

func (n *Node) StopVM(vmId uuid.UUID) error {
	return nil
}
