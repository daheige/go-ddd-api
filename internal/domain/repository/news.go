package repository

import (
	"github.com/daheige/go-ddd-api/internal/domain/model"
)

// NewsRepository represent repository of  the news
// Expect implementation by the infras layer
type NewsRepository interface {
	// Get obtain news by id
	Get(id int) (*model.News, error)
	// GetAll obtain all news
	GetAll() ([]model.News, error)
	// GetBySlug obtain news list by slug
	GetBySlug(slug string) ([]model.News, error)
	// GetAllByStatus obtain news by status
	GetAllByStatus(status string) ([]model.News, error)
	// Save news save
	Save(*model.News) error
	// Remove news remove by id
	Remove(id int) error
	// Update news update by entity
	Update(*model.News) error
}
