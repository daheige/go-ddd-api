package config

import (
	"log"

	"github.com/daheige/go-ddd-api/domain"
)

// DBMigrate will create & migrate the tables, then make the some relationships if necessary
func DBMigrate() error {
	err := AppConf.DB.AutoMigrate(domain.News{}, domain.Topic{}).Error
	log.Println("Migration error: ", err)
	log.Println("Migration has been processed")

	return err
}
