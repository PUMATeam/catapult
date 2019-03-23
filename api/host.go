package api

import (
	"net/http"

	"github.com/PUMATeam/catapult/model"
	"github.com/go-chi/render"
)

var host = new(model.Host)

// GetHosts retrieves all available hosts
func GetHosts(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, host.GetHosts())
}
