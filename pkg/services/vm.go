package services

import (
	"context"
	"fmt"

	"github.com/PUMATeam/catapult/pkg/node"

	"github.com/go-chi/chi/middleware"

	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/model"

	"github.com/PUMATeam/catapult/pkg/repositories"
	uuid "github.com/satori/go.uuid"
)

// NewVMsService instantiates a new VM service
func NewVMsService(vr repositories.VMs, hs Hosts, log *log.Logger) VMs {
	// TODO this looks weird and wrong
	vs := &vmsService{
		vmsRepository: vr,
		hostsService:  hs,
		log:           log,
	}

	return vs
}

type vmsService struct {
	vmsRepository repositories.VMs
	hostsService  Hosts
	log           *log.Logger
}

func (v *vmsService) AddVM(ctx context.Context, vm NewVM) (uuid.UUID, error) {
	// TODO extract to some util
	vmToAdd := model.VM{
		ID:             uuid.NewV4(),
		Name:           vm.Name,
		VCPU:           vm.VCPU,
		Memory:         vm.Memory,
		Status:         model.DOWN,
		HostID:         uuid.Nil,
		KernelImage:    vm.Kernel,
		RootFileSystem: vm.RootFileSystem,
	}

	v.log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
			"VM":        vmToAdd.Name,
		}).Info("adding VM")
	return v.vmsRepository.AddVM(ctx, vmToAdd)
}

func (v *vmsService) StartVM(ctx context.Context, vmID uuid.UUID) (*model.VM, error) {
	// TODO: algorithm should be - look for a host in status up and run the
	// VM on it
	h := v.findHostUP(ctx)
	if h == nil {
		return nil, fmt.Errorf("Could not find host in status up")
	}

	nodeService := node.NewNodeService(h, v.hostsService.GetConnManager(ctx))
	vm, err := v.VMByID(ctx, vmID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	cfg, err := nodeService.StartVM(ctx, vm)
	if err != nil {
		return nil, err
	}

	vm.HostID = h.ID
	vm.Address = cfg.GetAddress()
	vm.Status = model.UP
	v.vmsRepository.UpdateVM(ctx, vm)

	return &vm, nil
}

func (v *vmsService) ListVms(ctx context.Context) ([]model.VM, error) {
	return v.vmsRepository.ListVMs(ctx)
}

func (v *vmsService) StopVM(ctx context.Context, vm *model.VM) (uuid.UUID, error) {
	h, err := v.hostsService.HostByID(ctx, vm.HostID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Host %s not found", vm.HostID)
	}

	nodeService := node.NewNodeService(&h, v.hostsService.GetConnManager(ctx))
	err = nodeService.StopVM(ctx, vm.ID)
	if err != nil {
		return uuid.Nil, err
	}

	return vm.ID, nil
}

func (v *vmsService) ListVmsForHost(ctx context.Context, hostID uuid.UUID) ([]model.VM, error) {
	return nil, nil
}

func (v *vmsService) UpdateVMStatus(ctx context.Context, vm *model.VM, status model.Status) error {
	vm.Status = status
	return v.vmsRepository.UpdateVM(ctx, *vm)
}

func (v *vmsService) VMByID(ctx context.Context, vmID uuid.UUID) (model.VM, error) {
	vm, err := v.vmsRepository.VMByID(ctx, vmID)
	if err != nil {
		return model.VM{}, ErrNotFound
	}

	return vm, nil
}

func (v *vmsService) findHostUP(ctx context.Context) *model.Host {
	hosts, err := v.hostsService.ListHosts(ctx)
	v.log.WithContext(ctx).
		WithFields(log.Fields{
			"requestID": ctx.Value(middleware.RequestIDKey),
		}).Info("hosts found: ", hosts)
	if err != nil {
		log.Error(err)
	}

	for _, h := range hosts {
		if h.Status == model.UP {
			return &h
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
