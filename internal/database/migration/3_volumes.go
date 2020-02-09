package migration

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		log.Info("Ceating table volumes...")
		_, err := db.Exec(`CREATE TABLE volumes 
							(id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
							 status INT4,
							 description VARCHAR(255),
							 image VARCHAR(255) UNIQUE,
							 size BIGINT)`)
		return err
	}, func(db migrations.DB) error {
		log.Info("Dropping table volumes...")
		_, err := db.Exec(`DROP TABLE volumes`)
		return err
	})
}
