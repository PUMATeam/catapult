package node

import (
	"context"
	fmt "fmt"

	"github.com/go-chi/chi/middleware"

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
	Host        *model.Host
	connManager *Connections
}

// NewNodeService creates a Node instance
// TODO add logger
func NewNodeService(host *model.Host, connManager *Connections) NodeService {
	return &Node{
		Host:        host,
		connManager: connManager,
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
	f := func(conn *grpc.ClientConn) (interface{}, error) {
		client := NewNodeClient(conn)
		log.Info("Sending start vm request...")
		return client.StartVM(ctx, vmConfig)
	}

	conn := n.connManager.GetConnection(n.Host.ID)

	// This can happen if the node manager was restarted
	// manually, or there was an error
	var err error
	if conn == nil {
		address := fmt.Sprintf("%s:%d", n.Host.Address, n.Host.Port)
		conn, err = n.connManager.CreateConnection(ctx, n.Host.ID, address)

		if err != nil {
			return nil, err
		}
	}

	resp, err := runOnNode(conn, f)
	vmResp := resp.(*VmResponse)

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"vm":        vm.Name,
		}).Infof("Returned VM %v", vmResp.GetConfig())

	return vmResp.GetConfig(), err
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

	f := func(conn *grpc.ClientConn) (interface{}, error) {
		client := NewNodeClient(conn)
		return client.StopVM(ctx, uuid)
	}

	conn := n.connManager.GetConnection(n.Host.ID)
	_, err := runOnNode(conn, f)

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"vm":        vmID,
		}).Infof("Stopped VM %v", vmID)

	return err
}

type executeOnNode func(conn *grpc.ClientConn) (interface{}, error)

func runOnNode(conn *grpc.ClientConn, f executeOnNode) (interface{}, error) {
	resp, err := f(conn)
	return resp, err
}
