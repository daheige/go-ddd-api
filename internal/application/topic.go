package application

import (
	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/domain/repository"
)

// TopicService topic service
type TopicService struct {
	TopicRepo repository.TopicRepository `inject:""`
}

// GetTopic returns a topic by id
func (s *TopicService) GetTopic(id int) (*model.Topic, error) {
	return s.TopicRepo.Get(id)
}

// GetAllTopic return all topics
func (s *TopicService) GetAllTopic() ([]model.Topic, error) {
	return s.TopicRepo.GetAll()
}

// AddTopic saves new topic
func (s *TopicService) AddTopic(name string, slug string) error {
	u := &model.Topic{
		Name: name,
		Slug: slug,
	}
	return s.TopicRepo.Save(u)
}

// RemoveTopic do remove topic by id
func (s *TopicService) RemoveTopic(id int) error {
	return s.TopicRepo.Remove(id)
}

// UpdateTopic do update topic by id
func (s *TopicService) UpdateTopic(p model.Topic, id int) error {
	p.ID = uint(id)

	return s.TopicRepo.Update(&p)
}
