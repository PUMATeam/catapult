package model

import (
	"github.com/google/uuid"
	"net"
)

// Host represents a host in the system
type Host struct {
	ID      uuid.UUID `gorm:"primary_key"`
	Name    string    `json:"name"`
	Address net.IP    `json:"address"`
}

// GetHosts retrieves all the available hosts
func GetHosts() []Host {
	var hosts []Host
	db.Find(&hosts)

	return hosts
}
