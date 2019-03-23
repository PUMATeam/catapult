package model

import (
	"github.com/google/uuid"
)

// Host represents a host in the system
type Host struct {
	ID      uuid.UUID `json:"id" gorm:"primary_key"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}

// GetHosts retrieves all the available hosts
func (h Host) GetHosts() []Host {
	var hosts []Host
	db.Find(&hosts)

	return hosts
}
