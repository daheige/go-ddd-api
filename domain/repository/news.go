package repository

import (
	"github.com/daheige/go-ddd-api/domain/model"
)

// NewsRepository represent repository of  the news
// Expect implementation by the infrastructure layer
type NewsRepository interface {
	Get(id int) (*model.News, error)
	GetAll() ([]model.News, error)
	GetBySlug(slug string) ([]*model.News, error)
	GetAllByStatus(status string) ([]model.News, error)
	Save(*model.News) error
	Remove(id int) error
	Update(*model.News) error
}
