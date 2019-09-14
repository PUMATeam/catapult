package model

import (
	uuid "github.com/satori/go.uuid"
)

// Host represents the hosts table
type Host struct {
	ID       uuid.UUID `sql:"id,pk,type:uuid default gen_random_uuid()" json:"id"`
	Name     string    `sql:"name,type:varchar(50)" json:"name"`
	Address  string    `sql:"address,type:varchar(16)" json:"address"`
	Status   Status    `sql:"status,type:int4" json:"status"`
	User     string    `sql:"host_user,type:varchar(32)" json:"user"`
	Password string    `sql:"password,type:text" json:"password"`
}

// VM represents the vms table
type VM struct {
	ID             uuid.UUID `sql:"id,pk,type:uuid default gen_random_uuid()" json:"id"`
	Name           string    `sql:"name,type:varchar(50)" json:"name"`
	Status         Status    `sql:"status,type:int4" json:"status"`
	HostID         uuid.UUID `sql:"host_id,type:uuid" json:"host_id"`
	VCPU           int64     `sql:"vcpu,type:int4" json:"vcpu"`
	Memory         int64     `sql:"memory,type:int4" json:"memory"`
	KernelImage    string    `sql:"kernel,type:varchar(255)" json:"kernel"`
	RootFileSystem string    `sql:"root_file_system,type:varchar(255)" json:"rootfs"`
	Address        string    `sql:"address,type:varchar(16)" json:"address"`
}

type Status int

const (
	UNKNOWN Status = iota
	DOWN
	UP
	INSTALLING
	FAILED
)
