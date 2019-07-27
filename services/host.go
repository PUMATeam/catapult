package services

import (
	"context"
	"log"

	"github.com/PUMATeam/catapult/node"

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

// TODO: not working
func (as *hostsService) UpdateHostStatus(ctx context.Context, host *model.Host, status int) error {
	host.Status = status
	return as.hostsRepository.UpdateHost(ctx, host)
}

func (as *hostsService) AddHost(ctx context.Context, newHost NewHost) (uuid.UUID, error) {
	// TODO add validations and rollback in case ansible fails
	// also, run within a worker pool framework
	host := model.Host{
		Name:    newHost.Name,
		Address: newHost.Address,
		Status:  DOWN,
		User:    newHost.User,
		// TODO: encrypt the password
		Password: newHost.Password,
	}

	id, err := as.hostsRepository.AddHost(ctx, host)
	n := node.NewNodeService(host)
	err = n.InstallHost()
	if err != nil {
		log.Println("Updating status of host", host.Name, "to up")
		err = as.UpdateHostStatus(ctx, &host, UP)
		if err != nil {
			log.Println("Failed to update status of host", host.Name, "to up")
		}
	}

	return id, err
}

type NewHost struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// TODO: find a better way to do this
const (
	DOWN       int = 1
	INSTALLING int = 2
	UP         int = 3
)
