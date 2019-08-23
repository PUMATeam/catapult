package services

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/internal/database"
	"github.com/PUMATeam/catapult/internal/database/migration"
	"github.com/PUMATeam/catapult/pkg/repositories"
	"github.com/spf13/viper"
)

var repository repositories.Hosts

func initLog() {
	// TODO make configurable
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func setupTest(t *testing.T) {
	initLog()

	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")

	err := os.Chdir(dir)
	if err != nil {
		t.Error(err)
	}

	viper.SetDefault("db_config", "db.test.toml")
	err = migration.Migrate([]string{"init"})
	if err != nil {
		t.Log("err returned ", err)
	}

	err = migration.Migrate([]string{"up"})
	if err != nil {
		t.Error(err)
	}
	db, err := database.Connect()
	if err != nil {
		t.Error(err)
	}

	repository = repositories.NewHostsRepository(db)
}

func tearDown(t *testing.T) {
	t.Log("Tearing down...")
	migration.Reset()
}

func TestHostService(t *testing.T) {
	setupTest(t)
	defer tearDown(t)
	t.Run("test Add Host", testAddHost)
}

func testAddHost(t *testing.T) {
	svc := NewHostsService(repository)
	newHost := &NewHost{
		Name: "test_host",
		// TODO create the host as part of the setup
		Address:  "192.168.122.155",
		User:     "root",
		Password: "centos",
	}
	id, err := svc.AddHost(context.Background(), newHost)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Created host %s", id)
}
