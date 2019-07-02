package stores

import (
	"github.com/PUMATeam/catapult/model"
	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

// HostStore represents a host database accessor
type HostStore struct {
	db *pg.DB
}

// NewHostStore instantiates HostStore
func NewHostStore(db *pg.DB) *HostStore {
	return &HostStore{
		db: db,
	}
}

// Get gets a host by host ID.
func (s *HostStore) Get(id string) (*model.Host, error) {
	host := model.Host{ID: id}
	err := s.db.Select(&host)

	return &host, err
}

// Create creates a host
func (s *HostStore) Create(description string, address string) (*model.Host, error) {
	hostID := uuid.New().String()
	host := model.Host{
		ID:          hostID,
		Description: description,
		Address:     address,
	}

	err := s.db.Insert(&host)

	return &host, err
}
