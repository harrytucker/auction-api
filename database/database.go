package database

import (
	"github.com/jinzhu/gorm"
	// cannot edit generated main file so postgres dialect imported here
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// GetDB : get the database connection
func GetDB() (*gorm.DB, error) {
	log.Debug("Opening database connection")

	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to connect to database")
		return nil, err
	}

	return db, err
}
