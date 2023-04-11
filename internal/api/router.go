package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/daheige/go-ddd-api/internal/api/middleware"
	"github.com/daheige/go-ddd-api/internal/api/news"
	"github.com/daheige/go-ddd-api/internal/api/topics"
	"github.com/daheige/go-ddd-api/internal/infras/config"
	"github.com/daheige/go-ddd-api/internal/infras/migration"
	"github.com/daheige/go-ddd-api/internal/infras/utils"
	"github.com/gorilla/mux"
)

// RouterHandler api router
type RouterHandler struct {
	AppConfig     *config.AppConfig        `inject:""`
	TopicHandler  *topics.TopicHandler     `inject:""`
	NewsHandler   *news.NewsHandler        `inject:""`
	MigrateAction *migration.MigrateAction `inject:""`
}

// Router create router handler
func (s *RouterHandler) Router() *mux.Router {
	router := mux.NewRouter()

	router.StrictSlash(true)

	// install access log and recover handler
	router.Use(middleware.AccessLog, middleware.RecoverHandler)

	// not found handler
	router.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	// Index Route
	router.HandleFunc("/", s.home)
	router.HandleFunc("/api/v1", s.home)

	// News Route
	router.HandleFunc("/api/v1/news", s.NewsHandler.GetAllNews).Methods("GET")
	router.HandleFunc("/api/v1/news/{param}", s.NewsHandler.GetNews).Methods("GET")
	router.HandleFunc("/api/v1/news", s.NewsHandler.CreateNews).Methods("POST")
	router.HandleFunc("/api/v1/news/{news_id}", s.NewsHandler.RemoveNews).Methods("DELETE")
	router.HandleFunc("/api/v1/news/{news_id}", s.NewsHandler.UpdateNews).Methods("PUT")

	// Topic Route
	router.HandleFunc("/api/v1/topic", s.TopicHandler.GetAllTopic).Methods("GET")
	router.HandleFunc("/api/v1/topic/{topic_id}", s.TopicHandler.GetTopic).Methods("GET")
	router.HandleFunc("/api/v1/topic", s.TopicHandler.CreateTopic).Methods("POST")
	router.HandleFunc("/api/v1/topic/{topic_id}", s.TopicHandler.RemoveTopic).Methods("DELETE")
	router.HandleFunc("/api/v1/topic/{topic_id}", s.TopicHandler.UpdateTopic).Methods("PUT")

	// Migration Route
	router.HandleFunc("/api/v1/migrate", s.migrate)

	// router walk check
	if s.AppConfig.AppDebug {
		err := s.walk(router)
		if err != nil {
			fmt.Println("router walk error:", err)
		}
	}

	return router
}

func (s *RouterHandler) walk(router *mux.Router) error {
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}

		var queriesTemplates []string
		queriesTemplates, err = route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}

		var queriesRegexps []string
		queriesRegexps, err = route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}

		var methods []string
		methods, err = route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}

		return nil
	})

	if err != nil {
		fmt.Println("router walk error:", err)
	}

	return err
}

// migrate db migrate handler
func (s *RouterHandler) migrate(w http.ResponseWriter, r *http.Request) {
	err := s.MigrateAction.DBMigrate()
	if err != nil {
		log.Println("request_uri: ", r.RequestURI)
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	msg := "Success Migrate"
	utils.JSON(w, http.StatusOK, msg)
}

// home index handler
func (s *RouterHandler) home(w http.ResponseWriter, _ *http.Request) {
	utils.Respond(w, http.StatusOK, "GO DDD API")
}
