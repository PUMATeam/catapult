package repositories

import (
	"context"

	"github.com/PUMATeam/catapult/model"
	uuid "github.com/satori/go.uuid"
)

type Hosts interface {
	AddHost(context.Context, model.Host) (uuid.UUID, error)
	ListHosts(context.Context) ([]model.Host, error)
	HostByID(context.Context, uuid.UUID) (*model.Host, error)
	UpdateHost(context.Context, *model.Host) error
}

type VMs interface {
	AddVM(context.Context, model.VM) (uuid.UUID, error)
	ListVMs(context.Context) ([]model.VM, error)
	VMByID(context.Context, uuid.UUID) (*model.VM, error)
	UpdateVM(context.Context, model.VM) error
}
