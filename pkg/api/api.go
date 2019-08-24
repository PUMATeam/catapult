package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/middleware"

	"github.com/PUMATeam/catapult/pkg/repositories"
	"github.com/PUMATeam/catapult/pkg/services"
	"github.com/go-chi/chi"

	"github.com/PUMATeam/catapult/internal/database"
	uuid "github.com/satori/go.uuid"
)

var port int
var logger *log.Logger

func New(hs services.Hosts,
	vs services.VMs) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	hostsEndpoints(r, hs)
	vmsEndpoints(r, vs)

	return r
}

func Bootstrap(p int) http.Handler {
	port = p
	logger = InitLog()

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	hr := repositories.NewHostsRepository(db)
	hs := services.NewHostsService(hr, logger)

	vr := repositories.NewVMsRepository(db)
	vs := services.NewVMsService(vr, hr, logger)

	return New(hs, vs)
}

// Start start the server and listens on the provided port
func Start(h http.Handler) {
	server := http.Server{
		Handler: h,
		Addr:    ":" + strconv.Itoa(port),
	}

	initLog()
	installSignal()

	logger.Infof("Starting server, listening on: %v", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err)
	}
}

func installSignal() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		<-gracefulStop
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
