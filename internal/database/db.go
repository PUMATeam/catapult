package database

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

type config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

var conf config

// Connect connects to the database
func Connect() (*pg.DB, error) {
	// TODO maybe the config value should be passed instead?
	_, err := toml.DecodeFile(viper.GetString("db_config"), &conf)
	if err != nil {
		return nil, err
	}

	opts, err := pg.ParseURL(connectionURL(conf))
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)
	db.AddQueryHook(dbLogger{})

	return db, nil
}

func connectionURL(opts config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		opts.User, opts.Password, opts.Host, opts.Port, opts.Database)
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	s, _ := q.FormattedQuery()
	log.Debugf("Executed %v", s)
}
