package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// InitRoutes initializes the API routes
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/hosts", getHosts)

	return r
}
func getHosts(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "get hosts")
}
