package node

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"

	"github.com/PUMATeam/catapult/model"
	"github.com/PUMATeam/catapult/util"
	uuid "github.com/satori/go.uuid"
)

// NodeService exposes operations to perform on a host
type NodeService interface {
	InstallHost() error
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

// InstallHost installs prerequisits on the host
func (n *Node) InstallHost() error {
	hi := hostInstall{
		User:            n.host.User,
		FcVersion:       FcVersion,
		AnsiblePassword: n.host.Password,
	}
	ac := util.NewAnsibleCommand(util.SetupHostPlaybook,
		hi.User,
		n.host.Address,
		util.StructToMap(hi, strings.ToLower))
	err := ac.ExecuteAnsible()
	if err != nil {
		log.Println("Error during host install: ", err)
		return err
	}

	return nil
}

func (n *Node) StartVM(vm model.VM) error {
	// TODO: make configurable
	conn, err := grpc.Dial(fmt.Sprintf("%s:8888", n.host.Address), grpc.WithInsecure())
	if err != nil {
		return err
	}

	uuid := &UUID{
		Value: vm.ID.String(),
	}

	vmConfig := &VmConfig{
		VmID:   uuid,
		Memory: int32(vm.Memory),
		Vcpus:  int32(vm.VCPU),
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

// TODO make it configurable
const FcVersion = "0.15.0"

type hostInstall struct {
	User            string `json:"ignore"`
	AnsiblePassword string `json:"ansible_ssh_pass"`
	FcVersion       string
}
