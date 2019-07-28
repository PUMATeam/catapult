package services

import (
	"context"
	"log"
	"strings"

	"github.com/PUMATeam/catapult/util"

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

func (hs *hostsService) HostByID(ctx context.Context, id uuid.UUID) (*model.Host, error) {
	return hs.hostsRepository.HostByID(ctx, id)
}

func (hs *hostsService) ListHosts(ctx context.Context) ([]model.Host, error) {
	return hs.hostsRepository.ListHosts(ctx)
}

func (hs *hostsService) UpdateHostStatus(ctx context.Context, host model.Host, status int) error {
	host.Status = status
	return hs.hostsRepository.UpdateHost(ctx, host)
}

func (hs *hostsService) AddHost(ctx context.Context, newHost NewHost) (uuid.UUID, error) {
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

	id, err := hs.hostsRepository.AddHost(ctx, host)
	host.ID = id
	err = hs.InstallHost(host)
	if err != nil {
		return uuid.Nil, err
	}

	log.Println("Updating status of host", host.Name, "to up")
	err = hs.UpdateHostStatus(ctx, host, UP)
	if err != nil {
		log.Println("Failed to update status of host", host.Name, "to up")
	}

	return id, err
}

// InstallHost installs prerequisits on the host
// TODO: leaving it as public to allow a user add a host
// without installing right away
func (hs *hostsService) InstallHost(h model.Host) error {
	hi := hostInstall{
		User:            h.User,
		FcVersion:       FcVersion,
		AnsiblePassword: h.Password,
	}
	ac := util.NewAnsibleCommand(util.SetupHostPlaybook,
		h.User,
		h.Address,
		util.StructToMap(hi, strings.ToLower))
	err := ac.ExecuteAnsible()
	if err != nil {
		log.Println("Error during host install: ", err)
		return err
	}

	return nil
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

type hostInstall struct {
	User            string `json:"ignore"`
	AnsiblePassword string `json:"ansible_ssh_pass"`
	FcVersion       string
}

// TODO make it configurable
const FcVersion = "0.15.0"
