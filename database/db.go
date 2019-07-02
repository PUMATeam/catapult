package database

import (
	"fmt"

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
	_, err := toml.DecodeFile(viper.GetString("db_config"), &conf)
	if err != nil {
		return nil, err
	}

	opts, err := pg.ParseURL(connectionURL(conf))
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	return db, nil
}

func connectionURL(opts config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		opts.User, opts.Password, opts.Host, opts.Port, opts.Database)
}
