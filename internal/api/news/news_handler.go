package news

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/infras/utils"
	"github.com/gin-gonic/gin"
)

// NewsHandler news handler
type NewsHandler struct {
	NewsService *application.NewsService `inject:""`
}

// GetNews get news
func (s *NewsHandler) GetNews(c *gin.Context) {
	param := c.Param("param")
	// if param is numeric than search by news_id, otherwise
	// if alphabetic then search by topic.Slug
	newsID, err := strconv.Atoi(param)
	if err != nil {
		// param is alphabetic
		news, err2 := s.NewsService.GetNewsByTopic(param)
		if err2 != nil {
			utils.Error(c, http.StatusNotFound, err2.Error())
			return
		}

		utils.JSON(c, http.StatusOK, news)
		return
	}

	// param is numeric
	news, err := s.NewsService.GetNews(newsID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, news)
}

// GetAllNews get all news
func (s *NewsHandler) GetAllNews(c *gin.Context) {
	status := c.Query("status")
	// if status parameter exist draft|deleted|publish
	if status == "draft" || status == "deleted" || status == "publish" {
		news, err := s.NewsService.GetAllNewsByFilter(status)
		if err != nil {
			utils.Error(c, http.StatusNotFound, err.Error())
			return
		}

		utils.JSON(c, http.StatusOK, news)
		return
	}

	limit := c.Query("limit")
	page := c.Query("page")
	// if custom pagination exist news?limit=15&page=2
	if limit != "" && page != "" {
		limit, _ := strconv.Atoi(limit)
		page, _ := strconv.Atoi(page)
		if limit != 0 && page != 0 {
			news, err := s.NewsService.GetAllNews(limit, page)
			if err != nil {
				utils.Error(c, http.StatusNotFound, err.Error())
				return
			}

			utils.JSON(c, http.StatusOK, news)
			return
		}
	}

	news, err := s.NewsService.GetAllNews(15, 1) // 15, 1 default pagination
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, news)
}

// CreateNews create news
func (s *NewsHandler) CreateNews(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var p model.News
	if err := decoder.Decode(&p); err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
	}

	err := s.NewsService.AddNews(p)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusCreated, nil)
}

// RemoveNews remove news
func (s *NewsHandler) RemoveNews(c *gin.Context) {
	newsID, err := strconv.Atoi(c.Param("news_id"))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	err = s.NewsService.RemoveNews(newsID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, nil)
}

// UpdateNews update news
func (s *NewsHandler) UpdateNews(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var p model.News
	err := decoder.Decode(&p)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
	}

	newsID, err := strconv.Atoi(c.Param("news_id"))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	err = s.NewsService.UpdateNews(p, newsID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(c, http.StatusOK, nil)
}
