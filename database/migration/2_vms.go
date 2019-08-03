package migration

import (
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		log.Println("Ceating table vms...")
		_, err := db.Exec(`CREATE TABLE vms 
							(id UUID DEFAULT gen_random_uuid(),
							 name VARCHAR(50), 
							 status INT4,
							 host_id UUID REFERENCES hosts(id),
							 vcpu INTEGER,
							 memory INTEGER,
							 kernel VARCHAR(255),
							 root_file_system VARCHAR(255))`)
		return err
	}, func(db migrations.DB) error {
		log.Println("Dropping table vms...")
		_, err := db.Exec(`DROP TABLE vms`)
		return err
	})
}
