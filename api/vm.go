package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/PUMATeam/catapult/model"
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

	listVMsHandler := httptransport.NewServer(
		vmsEndpoint(vs),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)
	r.Method(http.MethodGet, "/vms", listVMsHandler)

	startVMHandler := httptransport.NewServer(
		startVMEndpoint(vs),
		decodeStartVMReq,
		encodeResponse,
	)
	r.Method(http.MethodGet, "/vms/{vmID}", startVMHandler)

}

func addVMEndpoint(svc services.VMs) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.NewVM)
		id, err := svc.AddVM(ctx, req)
		return IDResponse{ID: id}, err
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
		reqVM := request.(model.VM)
		vm, err := vs.StartVM(ctx, reqVM)
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

// TODO unify with decodeAddVmReq
func decodeStartVMReq(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var vm model.VM
	log.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&vm)
	return vm, err
}
