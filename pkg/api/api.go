package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/PUMATeam/catapult/pkg/storage"

	logrus "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/pkg/node"
	"github.com/PUMATeam/catapult/pkg/util"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/PUMATeam/catapult/pkg/repositories"
	"github.com/PUMATeam/catapult/pkg/services"
	"github.com/go-chi/chi"

	"github.com/PUMATeam/catapult/internal/database"
	uuid "github.com/satori/go.uuid"
)

var log *logrus.Logger
var connManager = node.NewNodeConnectionManager()

func newAPI(hs services.Hosts,
	vs services.VMs,
	vls services.Volumes) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(enableCors().Handler)
	hostsEndpoints(r, hs)
	vmsEndpoints(r, vs)
	volumesEndpoints(r, vls)

	return r
}

func bootstrap(log *logrus.Logger) http.Handler {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	hr := repositories.NewHostsRepository(db)

	hs := services.NewHostsService(hr, log, connManager)

	go func() {
		errors := hs.InitializeHosts(context.Background())
		if len(errors) > 0 {
			log.Error(errors)
		}

	}()

	vr := repositories.NewVMsRepository(db)
	vs := services.NewVMsService(vr, hs, log)
	ss := storage.NewStorageService(log)
	vrs := repositories.NewVolumesRepository(db)
	vls := services.NewVolumesService(hs, ss, vrs, log)

	return newAPI(hs, vs, vls)
}

// Start start the server and listens on the provided port
func Start(port int) {
	log = util.InitLog()
	handler := bootstrap(log)
	server := http.Server{
		Handler: handler,
		Addr:    ":" + strconv.Itoa(port),
	}

	installSignal()

	log.Infof("Starting server, listening on: %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}

func installSignal() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-gracefulStop
		log.Info("Shutting down grpc connections...")
		errors := connManager.Shutdown()
		if len(errors) > 0 {
			log.Warn(errors)
		}

		log.Info("Exiting... ")
		os.Exit(0)
	}()
}

type IDResponse struct {
	ID uuid.UUID `json:"id"`
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err,
	})
}

func codeFrom(err error) int {
	switch err {
	case services.ErrNotFound:
		return http.StatusNotFound
	case services.ErrAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func enableCors() *cors.Cors {
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	return cors
}
