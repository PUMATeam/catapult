package node

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/middleware"

	"google.golang.org/grpc/connectivity"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/PUMATeam/catapult/pkg/model"
	uuid "github.com/satori/go.uuid"
)

// NodeService exposes operations to perform on a host
type NodeService interface {
	ListVMs() ([]uuid.UUID, error)
	StartVM(ctx context.Context, vm model.VM) error
	StopVM(vmID uuid.UUID) error
}

type Node struct {
	Host *model.Host
}

// NewNodeService creates a Node instance
func NewNodeService(host *model.Host) NodeService {
	return &Node{
		Host: host,
	}
}

// StartVM starts an FC VM on the host
func (n *Node) StartVM(ctx context.Context, vm model.VM) error {
	// TODO: make port configurable
	conn, err := grpc.Dial(fmt.Sprintf("%s:8001", n.Host.Address),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	// Add a timeout and extract for reuse purposes
	for {
		if conn.GetState() == connectivity.Ready {
			log.Info("Connection is ready")
			break
		}
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

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"vm":        vm.Name,
		}).Infof("Returned VM %v", resp.GetConfig())

	return nil
}

// ListVMs lists VMs available on a node
func (n *Node) ListVMs() ([]uuid.UUID, error) {
	return nil, nil
}

// StopVM stops a running VM
func (n *Node) StopVM(vmID uuid.UUID) error {
	return nil
}
