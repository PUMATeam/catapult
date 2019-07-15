package migration

import (
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		log.Println("Adding encryption extension")
		_, err := db.Exec(`CREATE EXTENSION encryption`)
		return err
	}, func(db migrations.DB) error {
		log.Println("Dropping encryption extension")
		_, err := db.Exec(`DROP EXTENSION encryption`)
		return err
	})
}
