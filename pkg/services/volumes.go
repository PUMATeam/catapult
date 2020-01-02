package services

import (
	"context"
	"fmt"
	"math"

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
	path, size, err := nodeService.CreateDrive(ctx, volume.ImageName)
	if err != nil {
		logger.
			WithContext(ctx).
			Error(err)
		return uuid.Nil, err
	}

	// TODO: extract to util
	volSize := int64(math.Ceil(float64(size) / (1024 * 1024 * 1024)))
	logger.WithContext(ctx).Infof("Creating volume %v with size %s GiB",
		volID,
		volSize)
	_, err = v.storageService.Create(ctx, &storage.Volume{
		UUID: volID.String(),

		Size: volSize,
	})
	if err != nil {
		logger.
			WithContext(ctx).
			Error(err)
		return uuid.Nil, err
	}

	logger.WithContext(ctx).WithField("path", path).Infof("Created drive")

	nodeService.ConnectVolume(ctx,
		&node.Volume{VolumeID: volID.String(),
			PoolName: "volumes"})

	logger.WithContext(ctx).
		WithField("volume", volID).
		WithField("size", size).
		Infof("Mapped volume")

	return volID, err
}

type VolumeReq struct {
	ImageName string `json:"image"`
}
