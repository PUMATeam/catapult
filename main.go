package main

import (
	"net/http"

	"github.com/PUMATeam/catapult/api"
	"github.com/PUMATeam/catapult/model"
)

func main() {
	model.InitDB()
	router := api.InitRoutes()
	http.ListenAndServe(":3333", router)
}
