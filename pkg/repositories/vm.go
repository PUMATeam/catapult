package repositories

import (
	"context"
	"strings"

	"github.com/PUMATeam/catapult/pkg/model"
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

func (v *vmsRepository) AddVM(ctx context.Context, vm model.VM) (uuid.UUID, error) {
	err := v.db.WithContext(ctx).Insert(&vm)

	// TODO disgusting hack, find an actual solution
	if err != nil && strings.Contains(err.Error(), "cannot convert") {
		return vm.ID, nil
	}

	return vm.ID, err
}

func (v *vmsRepository) UpdateVM(ctx context.Context, vm model.VM) error {
	err := v.db.WithContext(ctx).Update(&vm)
	return err
}

func (v *vmsRepository) VMByID(ctx context.Context, vmID uuid.UUID) (model.VM, error) {
	vm := model.VM{
		ID: vmID,
	}

	err := v.db.WithContext(ctx).Select(&vm)

	// TODO disgusting hack, find an actual solution
	if err != nil && strings.Contains(err.Error(), "cannot convert") {
		return vm, nil
	}

	return vm, err
}

func (v *vmsRepository) ListVMs(ctx context.Context) ([]model.VM, error) {
	var vms []model.VM
	err := v.db.WithContext(ctx).Model(&vms).Select()
	return vms, err
}
