package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/domain/repository"
)

var _ repository.NewsRepository = (*NewsRepositoryImpl)(nil)

// NewsRepositoryImpl Implements repository.NewsRepository
type NewsRepositoryImpl struct {
	DB *gorm.DB `inject:""`
}

// Get news by id return domain.news
func (r *NewsRepositoryImpl) Get(id int) (*model.News, error) {
	news := &model.News{}
	if err := r.DB.Preload("Topic").First(&news, id).Error; err != nil {
		return nil, err
	}
	return news, nil
}

// GetAll News return all domain.news
func (r *NewsRepositoryImpl) GetAll() ([]model.News, error) {
	news := []model.News{}
	if err := r.DB.Preload("Topic").Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

// Save to add news
func (r *NewsRepositoryImpl) Save(news *model.News) error {
	if err := r.DB.Save(&news).Error; err != nil {
		return err
	}

	return nil
}

// Remove to delete news by id
func (r *NewsRepositoryImpl) Remove(id int) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	news := model.News{}
	if err := tx.First(&news, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	news.Status = "deleted"
	if err := tx.Save(&news).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&news).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Update is update news
func (r *NewsRepositoryImpl) Update(news *model.News) error {
	n := model.News{Title: news.Title, Slug: news.Slug,
		Content: news.Content, Status: news.Status, Topic: news.Topic}
	n.UpdatedAt = news.UpdatedAt
	if err := r.DB.Model(&news).UpdateColumns(n).Error; err != nil {
		return err
	}

	return nil
}

// GetAll News return all domain.news
func (r *NewsRepositoryImpl) GetAllByStatus(status string) ([]model.News, error) {
	if status == "deleted" {
		news := []model.News{}
		if err := r.DB.Unscoped().Where("status = ?", status).Preload("Topic").Find(&news).Error; err != nil {
			return nil, err
		}

		return news, nil
	}

	news := []model.News{}
	if err := r.DB.Where("status = ?", status).Preload("Topic").Find(&news).Error; err != nil {
		return nil, err
	}

	return news, nil
}

// GetBySlug News return all []model.News by topic.slug
func (r *NewsRepositoryImpl) GetBySlug(slug string) ([]model.News, error) {
	rows, err := r.DB.Raw("SELECT news.id, news.title, news.slug, news.content, news.status FROM `news_topics`"+
		" LEFT JOIN news ON news_topics.news_id=news.id WHERE "+
		"news_topics.topic_id=(SELECT id as topic_id FROM `topics`"+
		" WHERE slug = ?)", slug).Rows() // (*sql.Rows, error)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	us := make([]model.News, 0, 10)
	for rows.Next() {
		u := model.News{}
		err = rows.Scan(&u.ID, &u.Title, &u.Slug, &u.Content, &u.Status)

		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	return us, nil
}
