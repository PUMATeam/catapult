package api

import (
	"context"
	"github.com/PUMATeam/catapult/pkg/storage"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func storageEndpoints(r *chi.Mux) {
	createVolumeHandler := httptransport.NewServer(
		createVolumeEndpoint(),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)
	r.Method(http.MethodPost, "/storage", createVolumeHandler)
}

func createVolumeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		s := &storage.Storage{}
		resp, err := s.Create(ctx, &storage.Volume{UUID: "hello", Size: 1})
		return resp, err
	}
}
