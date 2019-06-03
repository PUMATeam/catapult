package main

import (
	"log"
	"net/http"

	"github.com/BurntSushi/toml"

	"github.com/PUMATeam/catapult/api"
	"github.com/PUMATeam/catapult/model"
)

func main() {
	if pgOptions, err := toml.DecodeFile("db.toml", &dbConfig); err != nil {
		log.Fatalf("Failed reading db config: %s", err)
	}

	db = model.Connect(pgOptions.database, pgOptions.username, pgOptions.password)
	defer db.Close()

	router := api.InitRoutes()
	http.ListenAndServe(":3333", router)
}
