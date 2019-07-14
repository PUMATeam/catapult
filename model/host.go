package model

import (
	uuid "github.com/satori/go.uuid"
)

// Host represents the hosts table
type Host struct {
	ID          uuid.UUID `json:"host_id"`
	Description string
	Address     string
	Status      int
}
