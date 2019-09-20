package migration

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		log.Info("Ceating table hosts...")
		_, err := db.Exec(`CREATE TABLE hosts 
							(id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
							 name VARCHAR(50) UNIQUE, 
							 address VARCHAR(16) UNIQUE,
							 port INT4,
							 status INT4,
							 host_user VARCHAR(32),
							 password TEXT)`)
		return err
	}, func(db migrations.DB) error {
		log.Info("Dropping table hosts...")
		_, err := db.Exec(`DROP TABLE hosts`)
		return err
	})
}
