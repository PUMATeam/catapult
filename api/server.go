package api

import (
	"log"
	"net/http"
	"strconv"
)

// Start start the server and listens on the provided port
func Start(port int) {
	api, err := New()
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Handler: api,
		Addr:    ":" + strconv.Itoa(port),
	}

	// TODO: add shutdown handling
	server.ListenAndServe()
}
