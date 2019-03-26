package api

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// InitRoutes initializes all middleware and routers
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/hosts", func(r chi.Router) {
		r.Get("/", GetHosts)
		r.Post("/", AddHost)
	})

	return r
}
