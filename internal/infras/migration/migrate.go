package migration

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/daheige/go-ddd-api/internal/domain/model"
)

// MigrateAction db migrate action
type MigrateAction struct {
	DB *gorm.DB `inject:""`
}

// DBMigrate will create & migrate the tables, then make the some relationships if necessary
func (m *MigrateAction) DBMigrate() error {
	err := m.DB.AutoMigrate(model.News{}, model.Topic{}).Error
	log.Println("Migration error: ", err)
	log.Println("Migration has been processed")

	return err
}
