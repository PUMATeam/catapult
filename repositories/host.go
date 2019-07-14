package repositories

import (
	"context"

	"github.com/PUMATeam/catapult/model"
	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
)

// NewHostsRepository create a new host repository
func NewHostsRepository(db *pg.DB) Hosts {
	return &hostsRepository{db: db}
}

type hostsRepository struct {
	db *pg.DB
}

func (h *hostsRepository) AddHost(ctx context.Context, host model.Host) (uuid.UUID, error) {
	return uuid.NewV4(), nil
}

func (h *hostsRepository) UpdateHost(ctx context.Context, host model.Host) error {
	return nil
}

func (h *hostsRepository) HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error) {
	return nil, nil
}

func (h *hostsRepository) ListHosts(ctx context.Context) ([]model.Host, error) {
	return nil, nil
}
