package host

import (
	"github.com/PUMATeam/catapult/database/stores"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
)

// API handles access to host related functions
type API struct {
	HostResource *Resource
}

// NewAPI instantiates a host endpoint API handler
func NewAPI(db *pg.DB) (*API, error) {
	// TODO: error handling
	hostStore := stores.NewHostStore(db)
	hostResource := NewHostResource(hostStore)
	api := &API{
		HostResource: hostResource,
	}

	return api, nil
}

// Router provides host routes
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/host", a.HostResource.router())
	return r
}
