package services

import (
	"context"
	"errors"

	"github.com/PUMATeam/catapult/model"
	uuid "github.com/satori/go.uuid"
)

type Hosts interface {
	HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error)
	ListHosts(ctx context.Context) ([]model.Host, error)
	AddHost(ctx context.Context, host NewHost) (uuid.UUID, error)
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)
