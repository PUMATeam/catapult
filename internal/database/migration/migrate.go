package migration

import (
	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/internal/database"
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
		log.Infof("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Infof("version is %d\n", oldVersion)
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
		log.Infof("migrated from version %d to %d\n", oldVersion, newVersion)
		version = newVersion
	}

	return nil
}
