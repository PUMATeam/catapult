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
	ListVMs(ctx context.Context) ([]uuid.UUID, error)
	StartVM(ctx context.Context, vm model.VM) error
	StopVM(ctx context.Context, vmID uuid.UUID) error
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

	// TODO: make port configurable
	conn, err := connectToNode(n.Host.Address, "8001")
	if err != nil {
		return err
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
func (n *Node) ListVMs(ctx context.Context) ([]uuid.UUID, error) {
	return nil, nil
}

// StopVM stops a running VM
func (n *Node) StopVM(ctx context.Context, vmID uuid.UUID) error {
	conn, err := connectToNode(n.Host.Address, "8001")
	if err != nil {
		return err
	}

	uuid := &UUID{
		Value: vmID.String(),
	}
	client := NewNodeClient(conn)
	resp, err := client.StopVM(ctx, uuid)

	log.WithContext(ctx).
	WithFields(log.Fields{
		"requestID": ctx.Value(middleware.RequestIDKey),
		"vm":        vmID,
	}).Infof("Stopped VM %v", resp.GetConfig())


	return err
}

func connectToNode(address string, port string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", address, port),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Add a timeout and extract for reuse purposes
	for {
		if conn.GetState() == connectivity.Ready {
			log.Info("Connection is ready")
			return conn, nil
		}
	}
}
