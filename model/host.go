package model

import (
	"github.com/jinzhu/gorm"
)

// Host represents a host in the system
type Host struct {
	gorm.Model
	ID      string `json:"id" gorm:"primary_key;type:uuid"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// GetHosts retrieves all the available hosts
func (h Host) GetHosts() []Host {
	var hosts []Host
	db.Find(&hosts)

	return hosts
}

// CreateHost adds the host to database
func (h Host) CreateHost(host Host) Host {
	if db.NewRecord(host) {
		db.Create(&host)
	}

	return host
}
