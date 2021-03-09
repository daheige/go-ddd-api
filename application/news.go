package application

import (
	"github.com/daheige/go-ddd-api/config"
	"github.com/daheige/go-ddd-api/domain"
	"github.com/daheige/go-ddd-api/infrastructure/pagination"
	"github.com/daheige/go-ddd-api/infrastructure/persistence"
)

// GetNews returns domain.news by id
func GetNews(id int) (*domain.News, error) {
	repo := persistence.NewNewsRepositoryWithRDB(config.AppConf.DB)
	return repo.Get(id)
}

// GetAllNews return all domain.news
func GetAllNews(limit int, page int) ([]domain.News, error) {
	var news []domain.News
	pagination.Paging(&pagination.Param{
		DB:      config.AppConf.DB.Preload("Topic"),
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
	}, &news)

	return news, nil
}

// AddNews saves new news
func AddNews(p domain.News) error {
	repo := persistence.NewNewsRepositoryWithRDB(config.AppConf.DB)
	return repo.Save(&p)
}

// RemoveNews do remove news by id
func RemoveNews(id int) error {
	repo := persistence.NewNewsRepositoryWithRDB(config.AppConf.DB)
	return repo.Remove(id)
}

// UpdateNews do remove news by id
func UpdateNews(p domain.News, id int) error {
	repo := persistence.NewNewsRepositoryWithRDB(config.AppConf.DB)
	p.ID = uint(id)

	return repo.Update(&p)
}

// GetAllNewsByFilter return all domain.news by filter
func GetAllNewsByFilter(status string) ([]domain.News, error) {
	repo := persistence.NewNewsRepositoryWithRDB(config.AppConf.DB)
	return repo.GetAllByStatus(status)
}

// GetNewsByTopic returns []domain.news by topic.slug
func GetNewsByTopic(slug string) ([]*domain.News, error) {
	repo := persistence.NewNewsRepositoryWithRDB(config.AppConf.DB)
	return repo.GetBySlug(slug)
}
