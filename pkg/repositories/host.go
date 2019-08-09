package repositories

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/model"
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
	host.ID = uuid.NewV4()
	err := h.db.WithContext(ctx).Insert(&host)
	return host.ID, err
}

func (h *hostsRepository) UpdateHost(ctx context.Context, host model.Host) error {
	err := h.db.WithContext(ctx).Update(&host)
	return err
}

func (h *hostsRepository) HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error) {
	host := &model.Host{ID: id}
	err := h.db.WithContext(ctx).Select(host)

	return host, err
}

func (h *hostsRepository) ListHosts(ctx context.Context) ([]model.Host, error) {
	var hosts []model.Host
	err := h.db.WithContext(ctx).Model(&hosts).Select()
	return hosts, err
}

func (h *hostsRepository) HostByAddress(ctx context.Context, address string) (*model.Host, error) {
	host := &model.Host{Address: address}
	err := h.db.WithContext(ctx).
		Model(host).
		WhereStruct(host).
		Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return host, nil
}

func (h *hostsRepository) HostByName(ctx context.Context, name string) (*model.Host, error) {
	host := &model.Host{Name: name}
	err := h.db.WithContext(ctx).
		Model(host).
		WhereStruct(host).
		Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return host, nil
}
