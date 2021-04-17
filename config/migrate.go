package config

import (
	"log"

	"github.com/daheige/go-ddd-api/domain/model"
)

// DBMigrate will create & migrate the tables, then make the some relationships if necessary
func DBMigrate() error {
	err := AppConf.DB.AutoMigrate(model.News{}, model.Topic{}).Error
	log.Println("Migration error: ", err)
	log.Println("Migration has been processed")

	return err
}
