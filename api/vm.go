package api

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/PUMATeam/catapult/services"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
)

func vmsEndpoints(r *chi.Mux, vs services.VMs) {
	addVMHandler := httptransport.NewServer(
		addVMEndpoint(vs),
		decodeAddVMReq,
		encodeResponse,
	)
	r.Method(http.MethodPost, "/vms", addVMHandler)

	startVMHandler := httptransport.NewServer(
		startVMEndpoint(vs),
		decodeVMByIDRequest,
		encodeResponse,
	)
	r.Method(http.MethodPost, "/vms/{vmID}/start", startVMHandler)

	listVMsHandler := httptransport.NewServer(
		vmsEndpoint(vs),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)
	r.Method(http.MethodGet, "/vms", listVMsHandler)
}

func addVMEndpoint(svc services.VMs) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.NewVM)
		id, err := svc.AddVM(ctx, req)
		log.Println("Returned from service.AddVM id", id)
		resp := IDResponse{ID: id}
		log.Println("Resp", resp)
		return resp, err
	}
}

func vmsEndpoint(vs services.VMs) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vms, err := vs.ListVms(ctx)
		return vms, err
	}
}

func startVMEndpoint(vs services.VMs) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqID := request.(uuid.UUID)
		vm, err := vs.StartVM(ctx, reqID)
		return vm, err
	}
}

func decodeAddVMReq(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var vm services.NewVM
	err := json.NewDecoder(r.Body).Decode(&vm)
	return vm, err
}

func decodeVMByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := uuid.FromString(chi.URLParam(r, "vmID"))
	if err != nil {
		return nil, err
	}
	return id, nil
}
