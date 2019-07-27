package services

import (
	"context"
	"log"

	"github.com/PUMATeam/catapult/node"

	"github.com/PUMATeam/catapult/model"
	"github.com/PUMATeam/catapult/repositories"
	uuid "github.com/satori/go.uuid"
)

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

// NewVMsService instantiates a new VM service
func NewVMsService(vr repositories.VMs, hr repositories.Hosts) VMs {
	// TODO this looks weird and wrong
	vs := &vmsService{
		vmsRepository:   vr,
		hostsRepository: hr,
	}

	// TODO: this should be done lazily
	vs.nodeService = vs.initNodeService()

	return vs
}

type vmsService struct {
	vmsRepository   repositories.VMs
	hostsRepository repositories.Hosts
	nodeService     node.NodeService
}

func (v *vmsService) AddVM(ctx context.Context, vm NewVM) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (v *vmsService) StartVM(ctx context.Context, vm model.VM) (*model.VM, error) {
	// TODO: algorithm should be - look for a host in status up and run the
	// VM on it
	err := v.nodeService.StartVM(vm)

	if err != nil {
		return nil, err
	}

	return &vm, nil
}

func (v *vmsService) ListVms(ctx context.Context) ([]uuid.UUID, error) {
	return nil, nil
}

func (v *vmsService) StopVM(ctx context.Context, host NewHost) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (v *vmsService) ListVmsForHost(ctx context.Context, hostID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

type NewVM struct {
	Name   string `json:"name"`
	VCPU   string `json:"vcpu"`
	Memory string `json:"memory"`
}
