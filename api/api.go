package api

import (
	"time"

	"github.com/PUMATeam/catapult/api/host"
	"github.com/PUMATeam/catapult/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New() (*chi.Mux, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	hostHandler, err := host.NewAPI(db)
	if err != nil {
		return nil, err
	}

	r := InitRoutes()
	r.Group(func(r chi.Router) {
		r.Mount("/", hostHandler.Router())
	})

	return r, nil
}

// InitRoutes initializes the API routes
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	return r
}
