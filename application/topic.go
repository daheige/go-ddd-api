package application

import (
	"github.com/daheige/go-ddd-api/config"
	"github.com/daheige/go-ddd-api/domain"
	"github.com/daheige/go-ddd-api/infrastructure/persistence"
)

// GetTopic returns a topic by id
func GetTopic(id int) (*domain.Topic, error) {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	return repo.Get(id)
}

// GetAllTopic return all topics
func GetAllTopic() ([]domain.Topic, error) {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	return repo.GetAll()
}

// AddTopic saves new topic
func AddTopic(name string, slug string) error {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	u := &domain.Topic{
		Name: name,
		Slug: slug,
	}
	return repo.Save(u)
}

// RemoveTopic do remove topic by id
func RemoveTopic(id int) error {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	return repo.Remove(id)
}

// UpdateTopic do update topic by id
func UpdateTopic(p domain.Topic, id int) error {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	p.ID = uint(id)

	return repo.Update(&p)
}
