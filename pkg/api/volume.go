package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PUMATeam/catapult/pkg/services"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func volumesEndpoints(r *chi.Mux, vls services.Volumes) {
	createVolumeHandler := httptransport.NewServer(
		addVolumeEndpoint(vls),
		decodeAddVolumesReq,
		encodeResponse,
	)
	r.Method(http.MethodPost, "/volumes", createVolumeHandler)
}

func addVolumeEndpoint(svc services.Volumes) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		volume := request.(services.VolumeReq)
		volID, err := svc.AddVolume(ctx, volume)
		return IDResponse{ID: volID}, err
	}
}

func decodeAddVolumesReq(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var volumeReq services.VolumeReq

	err := json.NewDecoder(r.Body).Decode(&volumeReq)
	if err != nil {
		return nil, err
	}

	return volumeReq, err
}
