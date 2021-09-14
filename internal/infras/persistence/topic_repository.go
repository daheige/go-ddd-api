package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/domain/repository"
)

var _ repository.TopicRepository = (*TopicRepositoryImpl)(nil)

// TopicRepositoryImpl Implements repository.TopicRepository
type TopicRepositoryImpl struct {
	DB *gorm.DB `inject:""`
}

// Get topic by id return domain.topic
func (r *TopicRepositoryImpl) Get(id int) (*model.Topic, error) {
	topic := &model.Topic{}
	if err := r.DB.Preload("News").First(&topic, id).Error; err != nil {
		return nil, err
	}
	return topic, nil
}

// GetAll topic return all domain.topic
func (r *TopicRepositoryImpl) GetAll() ([]model.Topic, error) {
	topics := []model.Topic{}
	if err := r.DB.Preload("News").Find(&topics).Error; err != nil {
		return nil, err
	}

	return topics, nil
}

// Save to add topic
func (r *TopicRepositoryImpl) Save(topic *model.Topic) error {
	if err := r.DB.Save(&topic).Error; err != nil {
		return err
	}

	return nil
}

// Remove delete topic
func (r *TopicRepositoryImpl) Remove(id int) error {
	topic := &model.Topic{}
	if err := r.DB.First(&topic, id).Error; err != nil {
		return err
	}

	if err := r.DB.Delete(&topic).Error; err != nil {
		return err
	}

	return nil
}

// Update data topic
func (r *TopicRepositoryImpl) Update(topic *model.Topic) error {
	if err := r.DB.Model(&topic).UpdateColumns(model.Topic{Name: topic.Name, Slug: topic.Slug}).Error; err != nil {
		return err
	}

	return nil
}
