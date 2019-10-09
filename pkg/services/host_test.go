package services

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"

	logrus "github.com/sirupsen/logrus"

	"github.com/PUMATeam/catapult/internal/database"
	"github.com/PUMATeam/catapult/internal/database/migration"
	"github.com/PUMATeam/catapult/pkg/repositories"
	"github.com/spf13/viper"
)

var repository repositories.Hosts
var log *logrus.Logger

func initLog() *logrus.Logger {
	// TODO make configurable
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	return log
}

func setupTest(t *testing.T) {
	log = initLog()

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
	svc := NewHostsService(repository, log)
	newHost := &NewHost{
		Name:     "test_host",
		Address:  "192.168.122.45",
		User:     "root",
		Password: "centos",
	}
	ctx := context.Background()
	id, err := svc.AddHost(ctx, newHost)
	if err != nil {
		t.Error(err)
	}

	h, err := svc.HostByID(ctx, id)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Created host %s", h.ID)
}
