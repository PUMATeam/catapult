package api

import (
	"fmt"
	"net/http"
)

// Start start the server and listens on the provided port
func Start(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), InitRoutes())
}
