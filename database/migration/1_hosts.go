package migration

import (
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		log.Println("Ceating table hosts...")
		_, err := db.Exec(`CREATE TABLE hosts 
							(id UUID DEFAULT gen_random_uuid(),
							 name VARCHAR(50), 
							 address VARCHAR(16),
							 status INT4)`)
		return err
	}, func(db migrations.DB) error {
		log.Println("Dropping table hosts...")
		_, err := db.Exec(`DROP TABLE hosts`)
		return err
	})
}
