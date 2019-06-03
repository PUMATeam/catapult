package model

import (
	"github.com/go-pg/pg"
)

// Connect connects to the database
func Connect(database string, username string, password string) (db *pg.DB) {
	return pg.Connect(&pg.Options{
		Database: database,
		User:     username,
		Password: password,
	})
}
