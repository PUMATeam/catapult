package services

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/model"

	"github.com/PUMATeam/catapult/pkg/node"

	"github.com/PUMATeam/catapult/pkg/repositories"
	uuid "github.com/satori/go.uuid"
)

// NewVMsService instantiates a new VM service
func NewVMsService(vr repositories.VMs, hr repositories.Hosts) VMs {
	// TODO this looks weird and wrong
	vs := &vmsService{
		vmsRepository:   vr,
		hostsRepository: hr,
	}

	return vs
}

type vmsService struct {
	vmsRepository   repositories.VMs
	hostsRepository repositories.Hosts
}

func (v *vmsService) AddVM(ctx context.Context, vm NewVM) (uuid.UUID, error) {
	// TODO extract to some util
	vmToAdd := model.VM{
		ID:             uuid.NewV4(),
		Name:           vm.Name,
		VCPU:           vm.VCPU,
		Memory:         vm.Memory,
		Status:         DOWN,
		HostID:         uuid.Nil,
		KernelImage:    vm.Kernel,
		RootFileSystem: vm.RootFileSystem,
	}

	log.Debug("add vm", vmToAdd)
	return v.vmsRepository.AddVM(ctx, vmToAdd)
}

func (v *vmsService) StartVM(ctx context.Context, vmID uuid.UUID) (*model.VM, error) {
	// TODO: algorithm should be - look for a host in status up and run the
	// VM on it
	nodeService := v.initNodeService()
	if nodeService == nil {
		return nil, fmt.Errorf("Could not find host in status up")
	}

	vm, err := v.VMByID(ctx, vmID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = nodeService.StartVM(ctx, vm)
	if err != nil {
		return nil, err
	}

	v.UpdateVMStatus(ctx, vm, UP)

	return &vm, nil
}

func (v *vmsService) ListVms(ctx context.Context) ([]model.VM, error) {
	return v.vmsRepository.ListVMs(ctx)
}

func (v *vmsService) StopVM(ctx context.Context, host NewHost) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (v *vmsService) ListVmsForHost(ctx context.Context, hostID uuid.UUID) ([]model.VM, error) {
	return nil, nil
}

func (v *vmsService) UpdateVMStatus(ctx context.Context, vm model.VM, status int) error {
	vm.Status = status
	return v.vmsRepository.UpdateVM(ctx, vm)
}

func (v *vmsService) VMByID(ctx context.Context, vmID uuid.UUID) (model.VM, error) {
	vm, err := v.vmsRepository.VMByID(ctx, vmID)
	if err != nil {
		return model.VM{}, ErrNotFound
	}

	return vm, nil
}

func (v *vmsService) initNodeService() node.NodeService {
	hosts, err := v.hostsRepository.ListHosts(context.TODO())
	log.Println("hosts found: ", hosts)
	if err != nil {
		log.Println(err)
	}

	for _, h := range hosts {
		if h.Status == UP {
			return node.NewNodeService(h)
		}
	}

	return nil
}

type NewVM struct {
	Name           string `json:"name"`
	VCPU           int64  `json:"vcpu"`
	Memory         int64  `json:"memory"`
	Kernel         string `json:"kernel"`
	RootFileSystem string `json:"rootfs"`
}
