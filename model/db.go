package model

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

var dbConfig config

var db *gorm.DB

func migrate() {
	db.AutoMigrate(&Host{})
}

// InitDB attempts to connect to the database
func InitDB() {
	if _, err := toml.DecodeFile("db.toml", &dbConfig); err != nil {
		log.Fatalf("Failed reading db config: %s", err)
	}

	var err error

	log.Print("Connecting to ", connection())
	db, err = gorm.Open("postgres", connection())
	if err != nil {
		log.Fatal(err)
	}

	migrate()
}

func connection() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
}
