package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/services"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
)

func hostsEndpoints(r *chi.Mux, hs services.Hosts) {
	addHostHandler := httptransport.NewServer(
		addHostEndpoint(hs),
		decodeAddHostReq,
		encodeResponse,
	)
	r.Method(http.MethodPost, "/hosts", addHostHandler)

	listHostsHandler := httptransport.NewServer(
		hostsEndpoint(hs),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)
	r.Method(http.MethodGet, "/hosts", listHostsHandler)

	hostByIDHandler := httptransport.NewServer(
		hostByIDEndpoint(hs),
		decodeHostByIDRequest,
		encodeResponse,
	)
	r.Method(http.MethodGet, "/hosts/{hostID}", hostByIDHandler)
}

func addHostEndpoint(svc services.Hosts) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.NewHost)
		if err := svc.Validate(ctx, req); err != nil {
			return nil, err
		}

		log.WithContext(ctx).
			WithFields(logrus.Fields{
				"requestID": ctx.Value(middleware.RequestIDKey),
				"request":   req,
			}).Info("Request to add host")
		id, err := svc.AddHost(ctx, &req)
		if err != nil {
			return nil, err
		}

		log.WithContext(ctx).
			WithFields(logrus.Fields{
				"requestID": ctx.Value(middleware.RequestIDKey),
				"host":      req.Name,
			}).Info("Added host")

		if req.ShouldInstall {
			h, err := svc.HostByID(ctx, id)
			if err != nil {
				return IDResponse{ID: id}, err
			}
			log.
				WithContext(ctx).
				WithFields(logrus.Fields{
					"requestID": ctx.Value(middleware.RequestIDKey),
					"host":      req.Name,
				}).Info("Installing host")
			go svc.InstallHost(ctx, &h, req.LocalNodePath)
		}

		return IDResponse{ID: id}, err
	}
}

func hostsEndpoint(hs services.Hosts) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		hosts, err := hs.ListHosts(ctx)
		return hosts, err
	}
}

func hostByIDEndpoint(svc services.Hosts) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uuid.UUID)
		host, err := svc.HostByID(ctx, req)
		return host, err
	}
}

func decodeAddHostReq(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var host services.NewHost
	err := json.NewDecoder(r.Body).Decode(&host)

	install := r.URL.Query().Get("install")
	if install != "" && install == "true" {
		host.ShouldInstall = true
	}

	return host, err
}

func decodeHostByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := uuid.FromString(chi.URLParam(r, "hostID"))
	if err != nil {
		return nil, err
	}
	return id, nil
}
