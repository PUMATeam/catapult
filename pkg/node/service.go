package node

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/middleware"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/PUMATeam/catapult/pkg/model"
	"github.com/PUMATeam/catapult/pkg/rpc"
	uuid "github.com/satori/go.uuid"
)

// NodeService exposes operations to perform on a host
type NodeService interface {
	ListVMs(ctx context.Context) ([]uuid.UUID, error)
	StartVM(ctx context.Context, vm model.VM) (*VmConfig, error)
	StopVM(ctx context.Context, vmID uuid.UUID) error
	CreateDrive(ctx context.Context, image string) (string, int64, error)
	ConnectVolume(ctx context.Context, volume *Volume) (string, error)
}

type Node struct {
	Host        *model.Host
	connManager *rpc.GRPCConnection
}

// NewNodeService creates a Node instance
// TODO add logger
func NewNodeService(host *model.Host, connManager *rpc.GRPCConnection) NodeService {
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

	address := fmt.Sprintf("%s:%d", n.Host.Address, n.Host.Port)

	f := func(conn *grpc.ClientConn) (interface{}, error) {
		client := NewNodeClient(conn)
		log.Info("Sending start vm request...")
		return client.StartVM(ctx, vmConfig)
	}

	conn := n.connManager.GetConnection(address)
	if conn == nil {
		return nil, fmt.Errorf("No connection found")
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

	conn := n.connManager.GetConnection(n.Host.Address)
	_, err := runOnNode(conn, f)

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"vm":        vmID,
		}).Infof("Stopped VM %v", vmID)

	return err
}

func (n *Node) CreateDrive(ctx context.Context, image string) (string, int64, error) {
	f := func(conn *grpc.ClientConn) (interface{}, error) {
		client := NewNodeClient(conn)
		return client.CreateDrive(ctx, &ImageName{Name: image})
	}

	conn := n.connManager.GetConnection(n.Host.Address)
	resp, err := runOnNode(conn, f)
	if err != nil {
		return "", -1, err
	}

	driveResp := resp.(*DriveResponse)

	return driveResp.GetPath(), driveResp.GetSize(), err
}

func (n *Node) ConnectVolume(ctx context.Context, volume *Volume) (string, error) {
	f := func(conn *grpc.ClientConn) (interface{}, error) {
		client := NewNodeClient(conn)
		return client.ConnectVolume(ctx, volume)
	}

	conn := n.connManager.GetConnection(n.Host.Address)
	resp, err := runOnNode(conn, f)
	connectResp := resp.(*ConnectResponse)
	if err != nil {
		return "", err
	}

	log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
		}).Infof("Connected volume %v", volume)

	return connectResp.GetPath(), nil
}

type executeOnNode func(conn *grpc.ClientConn) (interface{}, error)

// TODO properly handle when there is no network access to the host
func runOnNode(conn *grpc.ClientConn, f executeOnNode) (interface{}, error) {
	resp, err := f(conn)
	return resp, err
}
