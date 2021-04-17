package application

import (
	"github.com/daheige/go-ddd-api/config"
	"github.com/daheige/go-ddd-api/domain/model"
	"github.com/daheige/go-ddd-api/infrastructure/persistence"
)

// GetTopic returns a topic by id
func GetTopic(id int) (*model.Topic, error) {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	return repo.Get(id)
}

// GetAllTopic return all topics
func GetAllTopic() ([]model.Topic, error) {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	return repo.GetAll()
}

// AddTopic saves new topic
func AddTopic(name string, slug string) error {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	u := &model.Topic{
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
func UpdateTopic(p model.Topic, id int) error {
	repo := persistence.NewTopicRepositoryWithRDB(config.AppConf.DB)
	p.ID = uint(id)

	return repo.Update(&p)
}
