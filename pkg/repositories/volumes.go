package repositories

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/model"
	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
)

// NewVolumesRepository create a new vm repository
func NewVolumesRepository(db *pg.DB) Volumes {
	return &volumesRepository{db: db}
}

type volumesRepository struct {
	db *pg.DB
}

func (v *volumesRepository) AddVolume(ctx context.Context, volume model.Volume) (uuid.UUID, error) {
	tx, err := v.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	err = v.db.WithContext(ctx).Insert(&volume)

	// TODO disgusting hack, find an actual solution
	if err != nil && strings.Contains(err.Error(), "cannot convert") {
		return volume.ID, nil
	}
	if err != nil {
		log.Error(err)
		return uuid.Nil, tx.Rollback()
	}

	return volume.ID, err
}

func (v *volumesRepository) VolumeByID(ctx context.Context, volumeID uuid.UUID) (model.Volume, error) {
	volume := model.Volume{
		ID: volumeID,
	}

	err := v.db.WithContext(ctx).Select(&volume)

	// TODO disgusting hack, find an actual solution
	if err != nil && strings.Contains(err.Error(), "cannot convert") {
		return volume, nil
	}

	return volume, err
}

func (v *volumesRepository) ListVolumes(ctx context.Context) ([]model.Volume, error) {
	var volumes []model.Volume
	err := v.db.WithContext(ctx).Model(&volumes).Select()
	return volumes, err
}
