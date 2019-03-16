package model

import (
	"github.com/google/uuid"
	"net"
)

// Host represents a host in the system
type Host struct {
	ID      uuid.UUID `json:"host_id" gorm:"primary_key"`
	name    string    `json:"name"`
	address net.IP    `json:"address"`
}

// GetHosts retrieves all the available hosts
func GetHosts() []Host {
	var hosts []Host
	db.Find(&hosts)

	return hosts
}
