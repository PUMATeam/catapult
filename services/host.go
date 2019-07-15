package services

import (
	"context"

	"github.com/PUMATeam/catapult/model"
	"github.com/PUMATeam/catapult/repositories"
	uuid "github.com/satori/go.uuid"
)

func NewHostsService(hr repositories.Hosts) Hosts {
	return &hostsService{
		hostsRepository: hr,
	}
}

type hostsService struct {
	hostsRepository repositories.Hosts
}

func (as *hostsService) HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error) {
	return as.hostsRepository.HostByID(ctx, id)
}

func (as *hostsService) ListHosts(ctx context.Context) ([]model.Host, error) {
	return as.hostsRepository.ListHosts(ctx)
}

func (as *hostsService) AddHost(ctx context.Context, newHost NewHost) (uuid.UUID, error) {
	// TODO add validations
	host := model.Host{
		Name:    newHost.Name,
		Address: newHost.Address,
		Status:  DOWN,
	}

	id, err := as.hostsRepository.AddHost(ctx, host)

	return id, err
}

type NewHost struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

const (
	DOWN       int = 1
	INSTALLING int = 2
	UP         int = 3
)
