package model

// Host represents a host in the system
type Host struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// GetHosts retrieves all the available hosts
func (h Host) GetHosts() []Host {

	return hosts
}

// CreateHost adds the host to database
func (h Host) CreateHost(host Host) Host {

	return host
}
