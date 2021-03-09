package repository

import "github.com/daheige/go-ddd-api/domain"

// TopicRepository represent repository of the topic
// Expect implementation by the infrastructure layer
type TopicRepository interface {
	Get(id int) (*domain.Topic, error)
	GetAll() ([]domain.Topic, error)
	Save(*domain.Topic) error
	Remove(id int) error
	Update(*domain.Topic) error
}
