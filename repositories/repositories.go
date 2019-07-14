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
	UpdateHost(context.Context, model.Host) error
}
