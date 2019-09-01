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
var logger *logrus.Logger

func initLog() *logrus.Logger {
	// TODO make configurable
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	return logger
}

func setupTest(t *testing.T) {
	logger = initLog()

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
	svc := NewHostsService(repository, logger)
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
