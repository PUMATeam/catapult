package services

import (
	"context"
	"fmt"
	"strings"
	"sync"

	logrus "github.com/sirupsen/logrus"

	node "github.com/PUMATeam/catapult/pkg/node"
	"github.com/PUMATeam/catapult/pkg/util"

	"github.com/PUMATeam/catapult/pkg/model"
	"github.com/PUMATeam/catapult/pkg/repositories"
	uuid "github.com/satori/go.uuid"
)

func NewHostsService(hr repositories.Hosts, log *logrus.Logger, connManager *node.Connections) Hosts {
	return &hostsService{
		hostsRepository: hr,
		logger:          log,
		connManager:     connManager,
	}
}

type hostsService struct {
	hostsRepository repositories.Hosts
	logger          *logrus.Logger
	connManager     *node.Connections
}

// InitalizeHosts initializes running hosts when the app starts.
// Currently only creats grpc connection, soon will run health checks
func (hs *hostsService) InitializeHosts(ctx context.Context) []error {
	errors := make([]error, 0, 0)
	hosts, err := hs.hostsRepository.ListHosts(ctx)

	if err != nil {
		errors = append(errors, err)
		return errors
	}

	var wg sync.WaitGroup

	for _, host := range hosts {
		address := fmt.Sprintf("%s:%d", host.Address, host.Port)
		switch host.Status {
		case model.INSTALLING:
			hs.log(ctx, host.Name).Info("Host in status INSTALLING, moving to DOWN")
			hs.UpdateHostStatus(ctx, &host, model.DOWN)
		case model.UP:
			// Set host status to INITIALIZING during intialization
			hs.UpdateHostStatus(ctx, &host, model.INITIALIZING)
			go func(h *model.Host) {
				defer wg.Done()
				wg.Add(1)
				_, err := hs.connManager.CreateConnection(ctx, h.ID, address)

				if err == nil {
					hs.UpdateHostStatus(ctx, h, model.UP)
				} else {
					hs.log(ctx, h.Name).Error("Failed to initialize host connection")

					errors = append(errors, err)
					hs.UpdateHostStatus(ctx, h, model.DOWN)
				}
			}(&host)
		}
	}

	wg.Wait()

	return errors
}

func (hs *hostsService) HostByID(ctx context.Context, id uuid.UUID) (model.Host, error) {
	return hs.hostsRepository.HostByID(ctx, id)
}

func (hs *hostsService) ListHosts(ctx context.Context) ([]model.Host, error) {
	return hs.hostsRepository.ListHosts(ctx)
}

func (hs *hostsService) updateHostStatus(ctx context.Context, host *model.Host, status model.Status) error {
	host.Status = status
	return hs.hostsRepository.UpdateHost(ctx, host)
}

func (hs *hostsService) Validate(ctx context.Context, host NewHost) error {
	hs.log(ctx, host.Name).Info("Validating host")

	h, err := hs.hostsRepository.HostByAddress(ctx, host.Address)
	if err != nil {
		hs.log(ctx, h.Name).Error(err)
		return err
	}
	if h.ID != uuid.Nil && h.Status != model.FAILED {
		hs.log(ctx, host.Name).Errorf("Host with address %s already exists", host.Address)
		return ErrAlreadyExists
	}
	h, err = hs.hostsRepository.HostByName(ctx, host.Name)
	if err != nil {
		hs.log(ctx, "").Error(err)
		return err
	}
	if h.ID != uuid.Nil && h.Status != model.FAILED {
		hs.log(ctx, host.Name).Error("Host with this name already exists")

		return ErrAlreadyExists
	}

	return nil
}

func (hs *hostsService) AddHost(ctx context.Context, newHost *NewHost) (uuid.UUID, error) {
	host := model.Host{
		Name:    newHost.Name,
		Address: newHost.Address,
		Status:  model.DOWN,
		User:    newHost.User,
		// TODO: encrypt the password
		Password: newHost.Password,
		Port:     newHost.Port,
	}

	id, err := hs.hostsRepository.AddHost(ctx, &host)
	if err != nil {
		return uuid.Nil, err
	}

	host.ID = id
	return id, err
}

// InstallHost installs prerequisits on the host
// TODO: leaving it as public to allow a user add a host
// without installing right away
func (hs *hostsService) InstallHost(ctx context.Context, h *model.Host, localNodePath string) {
	hi := hostInstallData{
		User:            h.User,
		FcVersion:       fcVersion,
		AnsiblePassword: h.Password,
		LocalNodePath:   localNodePath,
		NodePort:        fmt.Sprintf("%d", h.Port),
	}

	hs.UpdateHostStatus(ctx, h, model.INSTALLING)

	ac := util.NewAnsibleCommand(util.SetupHostPlaybook,
		h.User,
		h.Address,
		util.StructToMap(hi, strings.ToLower),
		hs.logger)

	err := ac.ExecuteAnsible()
	if err != nil {
		hs.log(ctx, h.Name).Error("Error during host install: ", err)
		hs.UpdateHostStatus(ctx, h, model.FAILED)
		return
	}

	address := fmt.Sprintf("%s:%d", h.Address, h.Port)

	hs.log(ctx, h.Name).
		WithField("address", address).
		Info("Createting grpc connection for host...")
	_, err = hs.connManager.CreateConnection(ctx, h.ID, address)
	if err != nil {
		hs.log(ctx, h.Name).Error("Failed to create grpc connections, "+
			"will be retried upon sending a request: ", err)
	}

	// TODO send a health check to the host before
	hs.UpdateHostStatus(ctx, h, model.UP)
}

// activates catapult-node on the host
func (hs *hostsService) ActivateHost(ctx context.Context, h *model.Host) {
	hi := hostInstallData{
		User:            h.User,
		AnsiblePassword: h.Password,
		NodePort:        fmt.Sprintf("%d", h.Port),
	}

	hs.UpdateHostStatus(ctx, h, model.ACTIVATING)

	ac := util.NewAnsibleCommand(util.ActivateHostPlaybook,
		h.User,
		h.Address,
		util.StructToMap(hi, strings.ToLower),
		hs.logger)

	err := ac.ExecuteAnsible()
	if err != nil {
		hs.log(ctx, h.Name).Error("Error during catapult node activation: ", err)
		hs.UpdateHostStatus(ctx, h, model.FAILED)
		return
	}

	address := fmt.Sprintf("%s:%d", h.Address, h.Port)

	hs.log(ctx, h.Name).
		WithField("address", address).
		Info("Createting grpc connection for host...")
	_, err = hs.connManager.CreateConnection(ctx, h.ID, address)
	if err != nil {
		hs.log(ctx, h.Name).Error("Failed to create grpc connections, "+
			"will be retried upon sending a request: ", err)
	}

	// TODO send a health check to the host before
	hs.UpdateHostStatus(ctx, h, model.UP)
}

func (hs *hostsService) UpdateHostStatus(ctx context.Context, host *model.Host, status model.Status) error {
	hs.log(ctx, host.Name).Infof("Updating status of host to %d", status)

	err := hs.updateHostStatus(ctx, host, status)
	if err != nil {
		hs.log(ctx, host.Name).Errorf("Failed to update status of host to %d", status)
		return err
	}

	return nil
}

func (hs *hostsService) GetConnManager(ctx context.Context) *node.Connections {
	return hs.connManager
}

func (hs *hostsService) log(ctx context.Context, hostName string) *logrus.Entry {
	entry := util.Log(ctx, hs.logger)
	if hostName != "" {
		entry = entry.WithField("host", hostName)
	}
	return entry
}

type NewHost struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	User          string `json:"user"`
	Password      string `json:"password"`
	LocalNodePath string `json:"local_node_path"`
	ShouldInstall bool   `json:"-"`
	Port          int    `json:"port"`
}

type HostReinstall struct {
	ID            uuid.UUID `json:"id"`
	LocalNodePath string    `json:"local_node_path"`
}

type hostInstallData struct {
	User            string `json:"ignore"`
	AnsiblePassword string `json:"ansible_ssh_pass"`
	FcVersion       string
	LocalNodePath   string `json:"local_node_path"`
	NodePort        string `json:"node_port"`
}

// TODO make it configurable
const fcVersion = "0.20.0"
