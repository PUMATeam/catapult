package repositories

import (
	"context"
	"log"

	"github.com/PUMATeam/catapult/model"
	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
)

// NewVMsRepository create a new vm repository
func NewVMsRepository(db *pg.DB) VMs {
	return &vmsRepository{db: db}
}

type vmsRepository struct {
	db *pg.DB
}

func (v *vmsRepository) AddVM(context.Context, model.VM) (uuid.UUID, error) {
	log.Println("not implementd yet")
	return uuid.Nil, nil
}

func (v *vmsRepository) UpdateVM(context.Context, model.VM) error {
	log.Println("not implementd yet")
	return nil
}

func (v *vmsRepository) VMByID(context.Context, uuid.UUID) (*model.VM, error) {
	log.Println("not implementd yet")
	return nil, nil
}

func (v *vmsRepository) ListVMs(context.Context) ([]model.VM, error) {
	log.Println("not implementd yet")
	return nil, nil
}
