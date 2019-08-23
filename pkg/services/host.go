package services

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/util"

	"github.com/PUMATeam/catapult/pkg/model"
	"github.com/PUMATeam/catapult/pkg/repositories"
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

func (hs *hostsService) UpdateHostStatus(ctx context.Context, host model.Host, status model.Status) error {
	host.Status = status
	return hs.hostsRepository.UpdateHost(ctx, host)
}

func (hs *hostsService) Validate(ctx context.Context, host NewHost) error {
	log.Infof("Validating host %v", host)
	h, err := hs.hostsRepository.HostByAddress(ctx, host.Address)
	if err != nil {
		log.Error(err)
		return err
	}
	if h.ID != uuid.Nil {
		log.Errorf("Host with address %s already exists", host.Address)
		return ErrAlreadyExists
	}
	h, err = hs.hostsRepository.HostByName(ctx, host.Name)
	if err != nil {
		log.Error(err)
		return err
	}
	if h.ID != uuid.Nil {
		log.Errorf("Host with name %s already exists", host.Name)
		return ErrAlreadyExists
	}

	return nil
}

func (hs *hostsService) AddHost(ctx context.Context, newHost *NewHost) (uuid.UUID, error) {
	host := model.Host{
		Name:    newHost.Name,
		Address: newHost.Address,
		Status:  model.INSTALLING,
		User:    newHost.User,
		// TODO: encrypt the password
		Password: newHost.Password,
	}

	id, err := hs.hostsRepository.AddHost(ctx, host)
	if err != nil {
		return uuid.Nil, err
	}

	host.ID = id
	err = hs.InstallHost(host, newHost.LocalNodePath)
	if err != nil {
		return uuid.Nil, err
	}

	log.Infof("Updating status of host %v to up", host.Name)
	err = hs.UpdateHostStatus(ctx, host, model.UP)
	if err != nil {
		log.Errorf("Failed to update status of host %v to up", host.Name)
	}

	return id, err
}

// InstallHost installs prerequisits on the host
// TODO: leaving it as public to allow a user add a host
// without installing right away
func (hs *hostsService) InstallHost(h model.Host, localNodePath string) error {
	hi := hostInstall{
		User:            h.User,
		FcVersion:       fcVersion,
		AnsiblePassword: h.Password,
		LocalNodePath:   localNodePath,
	}
	log.Infof("Installing host %s", h.Name)
	ac := util.NewAnsibleCommand(util.SetupHostPlaybook,
		h.User,
		h.Address,
		util.StructToMap(hi, strings.ToLower))
	err := ac.ExecuteAnsible()
	if err != nil {
		log.Error("Error during host install: ", err)
		return err
	}

	return nil
}

type NewHost struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	User          string `json:"user"`
	Password      string `json:"password"`
	LocalNodePath string `json:"local_node_path"`
}

type hostInstall struct {
	User            string `json:"ignore"`
	AnsiblePassword string `json:"ansible_ssh_pass"`
	FcVersion       string
	LocalNodePath   string `json:"local_node_path"`
}

// TODO make it configurable
const fcVersion = "0.17.0"
