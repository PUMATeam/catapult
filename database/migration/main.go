package migration

import (
	"log"

	"github.com/PUMATeam/catapult/database"
	"github.com/go-pg/migrations"
)

// Migrate runs database migrations
func Migrate(args []string) error {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	oldVersion, newVersion, err := migrations.Run(db, args...)
	if err != nil {
		return err
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	return nil
}

// Reset resets the database
func Reset() error {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	version, err := migrations.Version(db)
	if err != nil {
		log.Fatal(err)
	}

	for version != 0 {
		oldVersion, newVersion, err := migrations.Run(db, "down")
		if err != nil {
			return err
		}
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		version = newVersion
	}

	return nil
}
