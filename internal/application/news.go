package application

import (
	"github.com/jinzhu/gorm"

	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/domain/repository"
	"github.com/daheige/go-ddd-api/internal/infras/pagination"
)

// NewsService news service
type NewsService struct {
	DB       *gorm.DB                  `inject:""`
	NewsRepo repository.NewsRepository `inject:""`
}

// GetNews returns domain.news by id
func (s *NewsService) GetNews(id int) (*model.News, error) {
	return s.NewsRepo.Get(id)
}

// GetAllNews return all domain.news
func (s *NewsService) GetAllNews(limit int, page int) ([]model.News, error) {
	var news []model.News
	pagination.Paging(&pagination.Param{
		DB:      s.DB.Preload("Topic"),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &news)

	return news, nil
}

// AddNews saves new news
func (s *NewsService) AddNews(p model.News) error {
	return s.NewsRepo.Save(&p)
}

// RemoveNews do remove news by id
func (s *NewsService) RemoveNews(id int) error {
	return s.NewsRepo.Remove(id)
}

// UpdateNews do remove news by id
func (s *NewsService) UpdateNews(p model.News, id int) error {
	p.ID = uint(id)

	return s.NewsRepo.Update(&p)
}

// GetAllNewsByFilter return all model.News by filter
func (s *NewsService) GetAllNewsByFilter(status string) ([]model.News, error) {
	return s.NewsRepo.GetAllByStatus(status)
}

// GetNewsByTopic returns []model.News by topic.slug
func (s *NewsService) GetNewsByTopic(slug string) ([]model.News, error) {
	return s.NewsRepo.GetBySlug(slug)
}
