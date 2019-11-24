package services

import (
	"context"
	"fmt"

	node "github.com/PUMATeam/catapult/pkg/node"
	"github.com/PUMATeam/catapult/pkg/repositories"
	"github.com/PUMATeam/catapult/pkg/storage"
	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
)

type volumesService struct {
	vmsRepository  repositories.VMs
	hostsService   Hosts
	storageService storage.Service
	logger         *logger.Logger
}

func NewVolumesService(hs Hosts, s storage.Service, logger *logger.Logger) Volumes {
	vls := &volumesService{
		hostsService:   hs,
		storageService: s,
		logger:         logger,
	}
	return vls
}

func (v *volumesService) AddVolume(ctx context.Context, volume VolumeReq) (uuid.UUID, error) {
	volID := uuid.NewV4()
	h := v.hostsService.FindHostUP(ctx)
	if h == nil {
		return uuid.Nil, fmt.Errorf("Could not find host in status up")
	}

	nodeService := node.NewNodeService(h, v.hostsService.GetConnManager(ctx))
	path, size, err := nodeService.CreateDrive(ctx, volume.imageName)
	v.storageService.Create(ctx, &storage.Volume{
		UUID: volID.String(),
		Size: size,
	})

	return volID, err
}

type VolumeReq struct {
	imageName string `json:"image"`
}
