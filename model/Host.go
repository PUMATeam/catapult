package model

// Host represents the hosts table
type Host struct {
	ID          string `json:"host_id"`
	Description string
	Address     string
}
