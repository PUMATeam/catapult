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

func (as *hostsService) AddHost(ctx context.Context) (uuid.UUID, error) {
	return uuid.NewV4(), nil
}
