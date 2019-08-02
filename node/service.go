package node

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/PUMATeam/catapult/model"
	uuid "github.com/satori/go.uuid"
)

// NodeService exposes operations to perform on a host
type NodeService interface {
	ListVMs() ([]uuid.UUID, error)
	StartVM(vm model.VM) error
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

func (n *Node) StartVM(vm model.VM) error {
	// TODO: make port configurable
	conn, err := grpc.Dial(fmt.Sprintf("%s:8888", n.host.Address), grpc.WithInsecure())
	if err != nil {
		return err
	}

	uuid := &UUID{
		Value: vm.ID.String(),
	}

	vmConfig := &VmConfig{
		VmID:   uuid,
		Memory: vm.Memory,
		Vcpus:  vm.VCPU,
	}
	client := NewNodeClient(conn)
	resp, err := client.StartVM(context.TODO(), vmConfig)
	if err != nil {
		fmt.Println("grpc error: ", err)
	}

	fmt.Println("grpc response:", resp)

	return nil
}

func (n *Node) ListVMs() ([]uuid.UUID, error) {
	return nil, nil
}

func (n *Node) StopVM(vmId uuid.UUID) error {
	return nil
}
