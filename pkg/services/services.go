package services

import (
	"context"
	"errors"

	"github.com/PUMATeam/catapult/pkg/model"
	uuid "github.com/satori/go.uuid"
)

type Hosts interface {
	Validate(ctx context.Context, host NewHost) error
	HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error)
	ListHosts(ctx context.Context) ([]model.Host, error)
	AddHost(ctx context.Context, host NewHost) (uuid.UUID, error)
	UpdateHostStatus(ctx context.Context, host model.Host, status model.Status) error
}

type VMs interface {
	AddVM(ctx context.Context, vm NewVM) (uuid.UUID, error)
	StartVM(ctx context.Context, id uuid.UUID) (*model.VM, error)
	ListVms(ctx context.Context) ([]model.VM, error)
	ListVmsForHost(ctx context.Context, hostID uuid.UUID) ([]model.VM, error)
	StopVM(ctx context.Context, host NewHost) (uuid.UUID, error)
	UpdateVMStatus(ctx context.Context, vm model.VM, status model.Status) error
	VMByID(ctx context.Context, vmID uuid.UUID) (model.VM, error)
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)
