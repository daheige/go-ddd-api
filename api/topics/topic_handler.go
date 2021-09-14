package topics

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/infras/utils"
)

// TopicHandler topic handler
type TopicHandler struct {
	TopicService *application.TopicService `inject:""`
}

// GetTopic get topic
func (s *TopicHandler) GetTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(mux.Vars(r)["topic_id"])
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	topic, err := s.TopicService.GetTopic(topicID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, topic)
}

// GetAllTopic get all topic
func (s *TopicHandler) GetAllTopic(w http.ResponseWriter, r *http.Request) {
	topics, err := s.TopicService.GetAllTopic()
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, topics)
}

type payload struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CreateTopic create topic
func (s *TopicHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var p payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	log.Println("p: ", p)
	err = s.TopicService.AddTopic(p.Name, p.Slug)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, nil)
}

// RemoveTopic remove topic
func (s *TopicHandler) RemoveTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(mux.Vars(r)["topic_id"])
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = s.TopicService.RemoveTopic(topicID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, nil)
}

// UpdateTopic update topic
func (s *TopicHandler) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p model.Topic
	err := decoder.Decode(&p)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
	}

	topicID, err := strconv.Atoi(mux.Vars(r)["topic_id"])
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = s.TopicService.UpdateTopic(p, topicID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, nil)
}
