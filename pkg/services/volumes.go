package services

import (
	"context"
	"fmt"
	"math"

	"github.com/PUMATeam/catapult/pkg/model"
	node "github.com/PUMATeam/catapult/pkg/node"
	"github.com/PUMATeam/catapult/pkg/repositories"
	"github.com/PUMATeam/catapult/pkg/storage"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type volumesService struct {
	volumesRepository repositories.Volumes
	hostsService      Hosts
	storageService    storage.Service
	log               *log.Logger
}

func NewVolumesService(hs Hosts, s storage.Service, vr repositories.Volumes, logger *log.Logger) Volumes {
	vls := &volumesService{
		hostsService:      hs,
		storageService:    s,
		volumesRepository: vr,
		log:               logger,
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
		v.log.WithContext(ctx).
			Error(err)
		return uuid.Nil, err
	}

	v.log.WithContext(ctx).
		WithFields(log.Fields{
			"volume": volID,
		}).
		Infof("Adding volume to database")
	_, err = v.volumesRepository.AddVolume(ctx, model.Volume{
		ID:          volID,
		Description: volume.Description,
		Status:      model.INITIALIZING,
		Image:       volume.ImageName,
		Size:        size,
	})
	if err != nil {
		v.log.WithContext(ctx).Error("Failed to add volume to database ", err)
		return uuid.Nil, err
	}
	// TODO: extract to util
	volSize := int64(math.Ceil(float64(size) / (1024 * 1024 * 1024)))
	v.log.WithContext(ctx).Infof("Creating volume %v with size %s GiB",
		volID,
		volSize)
	_, err = v.storageService.Create(ctx, &storage.Volume{
		UUID: volID.String(),
		Size: volSize,
	})
	if err != nil {
		v.log.WithContext(ctx).
			Error(err)
		return uuid.Nil, err
	}

	v.log.WithContext(ctx).WithField("path", path).Infof("Created drive")

	path, err = nodeService.ConnectVolume(ctx,
		&node.Volume{
			VolumeID:  volID.String(),
			PoolName:  "volumes",
			ImagePath: path})
	if err != nil {
		v.log.Errorf("Failed to connect volume %s", err)
	}

	v.log.WithContext(ctx).
		WithField("volume", volID).
		WithField("size", size).
		WithField("path", path).
		Infof("Mapped volume")

	return volID, err
}

type VolumeReq struct {
	ImageName   string `json:"image"`
	Description string `json:"description"`
}
