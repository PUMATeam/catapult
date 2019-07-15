package model

import (
	uuid "github.com/satori/go.uuid"
)

// Host represents the hosts table
type Host struct {
	ID      uuid.UUID `sql:"id,pk,type:uuid default gen_random_uuid()" json:"id"`
	Name    string
	Address string
	Status  int
}
