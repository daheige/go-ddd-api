package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/daheige/go-ddd-api/domain/model"
	"github.com/daheige/go-ddd-api/domain/repository"
)

// TopicRepositoryImpl Implements repository.TopicRepository
type TopicRepositoryImpl struct {
	Conn *gorm.DB
}

// NewTopicRepositoryWithRDB returns initialized TopicRepositoryImpl
func NewTopicRepositoryWithRDB(conn *gorm.DB) repository.TopicRepository {
	return &TopicRepositoryImpl{Conn: conn}
}

// Get topic by id return domain.topic
func (r *TopicRepositoryImpl) Get(id int) (*model.Topic, error) {
	topic := &model.Topic{}
	if err := r.Conn.Preload("News").First(&topic, id).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

// GetAll topic return all domain.topic
func (r *TopicRepositoryImpl) GetAll() ([]model.Topic, error) {
	topics := []model.Topic{}
	if err := r.Conn.Preload("News").Find(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}

// Save to add topic
func (r *TopicRepositoryImpl) Save(topic *model.Topic) error {
	if err := r.Conn.Save(&topic).Error; err != nil {
		return err
	}

	return nil
}

// Remove delete topic
func (r *TopicRepositoryImpl) Remove(id int) error {
	topic := &model.Topic{}
	if err := r.Conn.First(&topic, id).Error; err != nil {
		return err
	}

	if err := r.Conn.Delete(&topic).Error; err != nil {
		return err
	}

	return nil
}

// Update data topic
func (r *TopicRepositoryImpl) Update(topic *model.Topic) error {
	if err := r.Conn.Model(&topic).UpdateColumns(model.Topic{Name: topic.Name, Slug: topic.Slug}).Error; err != nil {
		return err
	}

	return nil
}
