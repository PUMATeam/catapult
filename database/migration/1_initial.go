package migration

import (
	"fmt"
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`CREATE EXTENSION encryption`)
		log.Println("creating table hosts...")
		_, err = db.Exec(`CREATE TABLE hosts 
							(host_id VARCHAR(36) PRIMARY KEY,
							 name VARCHAR(50), 
							 address VARCHAR(16))`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DROP EXTENSION encryption`)
		fmt.Println("dropping table hosts...")
		_, err = db.Exec(`DROP TABLE hosts`)
		return err
	})
}
