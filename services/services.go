package services

import (
	"context"
	"errors"

	"github.com/PUMATeam/catapult/node"

	"github.com/PUMATeam/catapult/model"
	uuid "github.com/satori/go.uuid"
)

type Hosts interface {
	HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error)
	ListHosts(ctx context.Context) ([]model.Host, error)
	AddHost(ctx context.Context, host NewHost) (uuid.UUID, error)
	UpdateHostStatus(ctx context.Context, host model.Host, status int) error
}

type VMs interface {
	AddVM(ctx context.Context, vm NewVM) (uuid.UUID, error)
	StartVM(ctx context.Context, vm node.RunVMCfg) (*node.RunVMCfg, error)
	ListVms(ctx context.Context) ([]uuid.UUID, error)
	ListVmsForHost(ctx context.Context, hostID uuid.UUID) ([]uuid.UUID, error)
	StopVM(ctx context.Context, host NewHost) (uuid.UUID, error)
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)
