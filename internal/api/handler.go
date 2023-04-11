package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/daheige/go-ddd-api/internal/api/middleware"
	"github.com/daheige/go-ddd-api/internal/api/news"
	"github.com/daheige/go-ddd-api/internal/api/topics"
	"github.com/daheige/go-ddd-api/internal/infras/config"
	"github.com/daheige/go-ddd-api/internal/infras/migration"
	"github.com/daheige/go-ddd-api/internal/infras/utils"
	"github.com/gorilla/mux"
)

var graceWait = 5 * time.Second

// NewsService application
type NewsService struct {
	AppConf       *config.AppConfig        `inject:""`
	TopicHandler  *topics.TopicHandler     `inject:""`
	NewsHandler   *news.NewsHandler        `inject:""`
	MigrateAction *migration.MigrateAction `inject:""`
}

// Run start services
func (s *NewsService) Run() {
	addr := fmt.Sprintf("0.0.0.0:%d", s.AppConf.Port)
	log.Printf("Server running on:%s", addr)

	// register mux router
	router := s.RouteHandler()

	// create http services
	server := &http.Server{
		// Handler: http.TimeoutHandler(router, time.Second*6, `{code:503,"message":"services timeout"}`),
		Handler:      router,
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// run http services in goroutine
	go func() {
		defer utils.Recover()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Println("services listen error:", err)
				return
			}

			log.Println("services will exit...")
		}
	}()

	// graceful exit
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// receive signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Block until we receive our signal.
	<-ch

	// Create s deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), graceWait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in s goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx)
	<-ctx.Done()

	log.Println("services shutdown success")
}

// RouteHandler returns the initialized router
func (s *NewsService) RouteHandler() *mux.Router {
	r := mux.NewRouter()

	r.StrictSlash(true)

	// install access log and recover handler
	r.Use(middleware.AccessLog, middleware.RecoverHandler)

	// not found handler
	r.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	// Index Route
	r.HandleFunc("/", s.home)
	r.HandleFunc("/api/v1", s.home)

	// News Route
	r.HandleFunc("/api/v1/news", s.NewsHandler.GetAllNews).Methods("GET")
	r.HandleFunc("/api/v1/news/{param}", s.NewsHandler.GetNews).Methods("GET")
	r.HandleFunc("/api/v1/news", s.NewsHandler.CreateNews).Methods("POST")
	r.HandleFunc("/api/v1/news/{news_id}", s.NewsHandler.RemoveNews).Methods("DELETE")
	r.HandleFunc("/api/v1/news/{news_id}", s.NewsHandler.UpdateNews).Methods("PUT")

	// Topic Route
	r.HandleFunc("/api/v1/topic", s.TopicHandler.GetAllTopic).Methods("GET")
	r.HandleFunc("/api/v1/topic/{topic_id}", s.TopicHandler.GetTopic).Methods("GET")
	r.HandleFunc("/api/v1/topic", s.TopicHandler.CreateTopic).Methods("POST")
	r.HandleFunc("/api/v1/topic/{topic_id}", s.TopicHandler.RemoveTopic).Methods("DELETE")
	r.HandleFunc("/api/v1/topic/{topic_id}", s.TopicHandler.UpdateTopic).Methods("PUT")

	// Migration Route
	r.HandleFunc("/api/v1/migrate", s.migrate)

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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

	return r
}

// migrate db migrate handler
func (s *NewsService) migrate(w http.ResponseWriter, r *http.Request) {
	err := s.MigrateAction.DBMigrate()
	if err != nil {
		utils.Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	msg := "Success Migrate"
	utils.JSON(w, http.StatusOK, msg)
}

// home index handler
func (s *NewsService) home(w http.ResponseWriter, _ *http.Request) {
	utils.Respond(w, http.StatusOK, "GO DDD API")
}
