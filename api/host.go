package api

import (
	"errors"
	"net/http"

	"github.com/PUMATeam/catapult/model"
	"github.com/go-chi/render"
)

var hostAccessor = new(model.Host)

type hostRequest struct {
	*model.Host
}

// GetHosts retrieves all available hosts
func GetHosts(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, hostAccessor.GetHosts())
}

// AddHost handles a request to create a new host
/* for example:
`{
	id: "9c604826-bce3-43ff-b2a3-8c20ccf3f9c7",
	name: "host",
	address: "192.168.1.1"
}`
*/
func AddHost(w http.ResponseWriter, r *http.Request) {
	data := &hostRequest{}

	if err := render.Bind(r, data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	host := data.Host
	hostAccessor.CreateHost(*host)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, host)
}

func (h *hostRequest) Bind(r *http.Request) error {
	if h.Host == nil {
		return errors.New("missing required host fields")
	}

	return nil
}
