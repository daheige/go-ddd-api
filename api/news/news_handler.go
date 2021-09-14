package news

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/daheige/go-ddd-api/internal/application"
	"github.com/daheige/go-ddd-api/internal/domain/model"
	"github.com/daheige/go-ddd-api/internal/infras/utils"
)

// NewsHandler news handler
type NewsHandler struct {
	NewsService *application.NewsService `inject:""`
}

// GetNews get news
func (s *NewsHandler) GetNews(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["param"]

	// if param is numeric than search by news_id, otherwise
	// if alphabetic then search by topic.Slug
	newsID, err := strconv.Atoi(param)
	if err != nil {
		// param is alphabetic
		news, err2 := s.NewsService.GetNewsByTopic(param)
		if err2 != nil {
			utils.Error(w, http.StatusNotFound, err2, err2.Error())
			return
		}

		utils.JSON(w, http.StatusOK, news)
		return
	}

	// param is numeric
	news, err := s.NewsService.GetNews(newsID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, news)
}

// GetAllNews get all news
func (s *NewsHandler) GetAllNews(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	status := queryValues.Get("status")

	// if status parameter exist draft|deleted|publish
	if status == "draft" || status == "deleted" || status == "publish" {
		news, err := s.NewsService.GetAllNewsByFilter(status)
		if err != nil {
			utils.Error(w, http.StatusNotFound, err, err.Error())
			return
		}

		utils.JSON(w, http.StatusOK, news)
		return
	}

	limit := queryValues.Get("limit")
	page := queryValues.Get("page")

	// if custom pagination exist news?limit=15&page=2
	if limit != "" && page != "" {
		limit, _ := strconv.Atoi(limit)
		page, _ := strconv.Atoi(page)

		if limit != 0 && page != 0 {
			news, err := s.NewsService.GetAllNews(limit, page)
			if err != nil {
				utils.Error(w, http.StatusNotFound, err, err.Error())
				return
			}

			utils.JSON(w, http.StatusOK, news)
			return
		}
	}

	news, err := s.NewsService.GetAllNews(15, 1) // 15, 1 default pagination
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, news)
}

// CreateNews create news
func (s *NewsHandler) CreateNews(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p model.News
	if err := decoder.Decode(&p); err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
	}

	err := s.NewsService.AddNews(p)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, nil)
}

// RemoveNews remove news
func (s *NewsHandler) RemoveNews(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(mux.Vars(r)["news_id"])
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = s.NewsService.RemoveNews(newsID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, nil)
}

// UpdateNews update news
func (s *NewsHandler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p model.News
	err := decoder.Decode(&p)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
	}

	newsID, err := strconv.Atoi(mux.Vars(r)["news_id"])
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = s.NewsService.UpdateNews(p, newsID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, nil)
}
