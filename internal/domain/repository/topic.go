package repository

import (
	"github.com/daheige/go-ddd-api/internal/domain/model"
)

// TopicRepository represent repository of the topic
// Expect implementation by the infras layer
type TopicRepository interface {
	Get(id int) (*model.Topic, error)
	GetAll() ([]model.Topic, error)
	Save(*model.Topic) error
	Remove(id int) error
	Update(*model.Topic) error
}
