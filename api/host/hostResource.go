package host

import (
	"net/http"

	"github.com/PUMATeam/catapult/database/stores"
	"github.com/go-chi/chi"
)

// Resource accesses database store for hosts
type Resource struct {
	Store *stores.HostStore
}

// NewHostResource creates a HostResource
func NewHostResource(store *stores.HostStore) *Resource {
	return &Resource{
		Store: store,
	}
}

func (hr *Resource) router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hosts root"))
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create a new Host"))
	})

	r.Route("/{hostID}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Get host by ID"))
		})
	})

	return r
}
