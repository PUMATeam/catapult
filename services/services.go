package services

import (
	"context"

	"github.com/PUMATeam/catapult/model"
	uuid "github.com/satori/go.uuid"
)

type Hosts interface {
	HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error)
	ListHosts(ctx context.Context) ([]model.Host, error)
	AddHost(ctx context.Context) (uuid.UUID, error)
}
