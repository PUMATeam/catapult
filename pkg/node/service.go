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
	StartVM(ctx context.Context, vm model.VM) (*VmConfig, error)
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
func (n *Node) StartVM(ctx context.Context, vm model.VM) (*VmConfig, error) {
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
	f := func(conn *grpc.ClientConn) (*Response, error) {
		client := NewNodeClient(conn)
		return client.StartVM(ctx, vmConfig)
	}

	resp, err := runOnNode(n.Host, f)

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"vm":        vm.Name,
		}).Infof("Returned VM %v", resp.GetConfig())

	return resp.GetConfig(), err
}

// ListVMs lists VMs available on a node
func (n *Node) ListVMs(ctx context.Context) ([]uuid.UUID, error) {
	return nil, nil
}

// StopVM stops a running VM
func (n *Node) StopVM(ctx context.Context, vmID uuid.UUID) error {
	uuid := &UUID{
		Value: vmID.String(),
	}

	f := func(conn *grpc.ClientConn) (*Response, error) {
		client := NewNodeClient(conn)
		return client.StopVM(ctx, uuid)
	}
	resp, err := runOnNode(n.Host, f)

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"vm":        vmID,
		}).Infof("Stopped VM %v", resp.GetConfig())

	return err
}

type executeOnNode func(conn *grpc.ClientConn) (*Response, error)

func runOnNode(host *model.Host, f executeOnNode) (*Response, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host.Address, host.Port),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	// Add a timeout and extract for reuse purposes
	for {
		if conn.GetState() == connectivity.Ready {
			log.Info("Connection is ready")
			resp, err := f(conn)
			return resp, err
		}
	}
}
