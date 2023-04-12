package topics

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/infras/utils"
	"github.com/gin-gonic/gin"
)

// TopicHandler topic handler
type TopicHandler struct {
	TopicService *application.TopicService `inject:""`
}

// GetTopic get topic
func (s *TopicHandler) GetTopic(c *gin.Context) {
	topicID, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	topic, err := s.TopicService.GetTopic(topicID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, topic)
}

// GetAllTopic get all topic
func (s *TopicHandler) GetAllTopic(c *gin.Context) {
	topics, err := s.TopicService.GetAllTopic()
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, topics)
}

type payload struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CreateTopic create topic
func (s *TopicHandler) CreateTopic(c *gin.Context) {
	var p payload
	err := json.NewDecoder(c.Request.Body).Decode(&p)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	err = s.TopicService.AddTopic(p.Name, p.Slug)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusCreated, nil)
}

// RemoveTopic remove topic
func (s *TopicHandler) RemoveTopic(c *gin.Context) {
	topicID, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	err = s.TopicService.RemoveTopic(topicID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, nil)
}

// UpdateTopic update topic
func (s *TopicHandler) UpdateTopic(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var topic model.Topic
	err := decoder.Decode(&topic)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
	}

	topicID, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	err = s.TopicService.UpdateTopic(topic, topicID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, nil)
}
