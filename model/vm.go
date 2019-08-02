package model

import (
	uuid "github.com/satori/go.uuid"
)

// VM represents the vms table
type VM struct {
	ID     uuid.UUID `sql:"id,pk,type:uuid default gen_random_uuid()" json:"id"`
	Name   string    `sql:"name,type:varchar(50)" json:"name"`
	Status int       `sql:"status,type:int4" json:"status"`
	HostID uuid.UUID `sql:"host_id,type:uuid" json:"host_id"`
	VCPU   int64     `sql:"vcpu,type:int4" json:"vcpu"`
	Memory int64     `sql:"memory,type:int4" json:"memory"`
}
