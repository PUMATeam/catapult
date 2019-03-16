package model

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	"log"
)

type config struct {
	Host     string
	Name     string
	User     string
	Password string
}

var dbConfig config

var db *gorm.DB

// InitDB attempts to connect to the database
func InitDB() {
	if _, err := toml.Decode("db.toml", &dbConfig); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(connection())
	if err != nil {
		log.Fatal(err)
	}
}

func connection() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name)
}

func Migrate() {

}
