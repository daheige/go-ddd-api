package interfaces

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/daheige/go-ddd-api/application"
	"github.com/daheige/go-ddd-api/config"
	"github.com/daheige/go-ddd-api/domain"
	"github.com/daheige/go-ddd-api/infrastructure/utils"
	"github.com/daheige/tigago/gutils"
	"github.com/gorilla/mux"
)

// IsLetter function to check string is aplhanumeric only
var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

var graceWait = 5 * time.Second

// Run start server
func Run(port int) {
	log.Printf("Server running on port:%d/", port)

	// register mux router
	router := RouteHandler()

	// create http server
	server := &http.Server{
		// Handler: http.TimeoutHandler(router, time.Second*6, `{code:503,"message":"server timeout"}`),
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// run http server in goroutine
	go func() {
		defer utils.Recover()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Println("server listen error:", err)
				return
			}

			log.Println("server will exit...")
		}
	}()

	// graceful exit
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recv signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), graceWait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx)
	<-ctx.Done()

	log.Println("server shutdown success")
}

// RouteHandler returns the initialized router
func RouteHandler() *mux.Router {
	r := mux.NewRouter()

	r.StrictSlash(true)

	// install access log and recover handler
	r.Use(AccessLog, RecoverHandler)

	// not found handler
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// Index Route
	r.HandleFunc("/", index)
	r.HandleFunc("/api/v1", index)

	// News Route
	r.HandleFunc("/api/v1/news", getAllNews)
	r.HandleFunc("/api/v1/news/{param}", getNews)
	r.HandleFunc("/api/v1/news", createNews)
	r.HandleFunc("/api/v1/news/{news_id}", removeNews).Methods("DELETE")
	r.HandleFunc("/api/v1/news/{news_id}", updateNews).Methods("PUT")

	// Topic Route
	r.HandleFunc("/api/v1/topic", getAllTopic)
	r.HandleFunc("/api/v1/topic/{topic_id}", getTopic)
	r.HandleFunc("/api/v1/topic", createTopic).Methods("POST")
	r.HandleFunc("/api/v1/topic/{topic_id}", removeTopic).Methods("DELETE")
	r.HandleFunc("/api/v1/topic/{topic_id}", updateTopic).Methods("PUT")

	// Migration Route
	r.HandleFunc("/api/v1/migrate", migrate)

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

// NotFoundHandler not found api router
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("this page not found"))
}

// RecoverHandler recover handler
func RecoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("exec panic error", map[string]interface{}{
					"trace_error": string(debug.Stack()),
				})

				// server error
				http.Error(w, "server error!", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func AccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		log.Println("exec begin", nil)
		log.Println("request before")
		log.Println("request uri: ", r.RequestURI)

		// x-request-id
		reqId := r.Header.Get("x-request-id")
		if reqId == "" {
			reqId = gutils.Uuid()
		}

		// log.Println("log_id: ", reqId)
		r = utils.ContextSet(r, "log_id", reqId)
		r = utils.ContextSet(r, "client_ip", r.RemoteAddr)
		r = utils.ContextSet(r, "request_method", r.Method)
		r = utils.ContextSet(r, "request_uri", r.RequestURI)
		r = utils.ContextSet(r, "user_agent", r.Header.Get("User-Agent"))

		h.ServeHTTP(w, r)

		log.Println("exec end", map[string]interface{}{
			"exec_time": time.Now().Sub(t).Seconds(),
		})

	})
}

func index(w http.ResponseWriter, _ *http.Request) {
	Respond(w, http.StatusOK, "GO DDD API")
}

// =============================
//    NEWS
// =============================

func getNews(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["param"]

	// if param is numeric than search by news_id, otherwise
	// if alphabetic then search by topic.Slug
	newsID, err := strconv.Atoi(param)
	if err != nil {
		// param is alphabetic
		news, err2 := application.GetNewsByTopic(param)
		if err2 != nil {
			Error(w, http.StatusNotFound, err2, err2.Error())
			return
		}

		JSON(w, http.StatusOK, news)
		return
	}

	// param is numeric
	news, err := application.GetNews(newsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, news)
}

func getAllNews(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	status := queryValues.Get("status")

	// if status parameter exist draft|deleted|publish
	if status == "draft" || status == "deleted" || status == "publish" {
		news, err := application.GetAllNewsByFilter(status)
		if err != nil {
			Error(w, http.StatusNotFound, err, err.Error())
			return
		}

		JSON(w, http.StatusOK, news)
		return
	}

	limit := queryValues.Get("limit")
	page := queryValues.Get("page")

	// if custom pagination exist news?limit=15&page=2
	if limit != "" && page != "" {
		limit, _ := strconv.Atoi(limit)
		page, _ := strconv.Atoi(page)

		if limit != 0 && page != 0 {
			news, err := application.GetAllNews(limit, page)
			if err != nil {
				Error(w, http.StatusNotFound, err, err.Error())
				return
			}

			JSON(w, http.StatusOK, news)
			return
		}
	}

	news, err := application.GetAllNews(15, 1) // 15, 1 default pagination
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, news)
}

func createNews(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p domain.News
	if err := decoder.Decode(&p); err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
	}

	err := application.AddNews(p)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusCreated, nil)
}

func removeNews(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(mux.Vars(r)["news_id"])
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.RemoveNews(newsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

func updateNews(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p domain.News
	err := decoder.Decode(&p)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
	}

	newsID, err := strconv.Atoi(mux.Vars(r)["news_id"])
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.UpdateNews(p, newsID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

// =============================
//    TOPIC
// =============================

func getTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(mux.Vars(r)["topic_id"])
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	topic, err := application.GetTopic(topicID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, topic)
}

func getAllTopic(w http.ResponseWriter, r *http.Request) {
	topics, err := application.GetAllTopic()
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, topics)
}

func createTopic(w http.ResponseWriter, r *http.Request) {

	type payload struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	var p payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.AddTopic(p.Name, p.Slug)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusCreated, nil)
}

func removeTopic(w http.ResponseWriter, r *http.Request) {
	topicID, err := strconv.Atoi(mux.Vars(r)["topic_id"])
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.RemoveTopic(topicID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

func updateTopic(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p domain.Topic
	err := decoder.Decode(&p)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
	}

	topicID, err := strconv.Atoi(mux.Vars(r)["topic_id"])
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	err = application.UpdateTopic(p, topicID)
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	JSON(w, http.StatusOK, nil)
}

// =============================
//    MIGRATE
// =============================

func migrate(w http.ResponseWriter, r *http.Request) {
	err := config.DBMigrate()
	if err != nil {
		Error(w, http.StatusNotFound, err, err.Error())
		return
	}

	msg := "Success Migrate"
	JSON(w, http.StatusOK, msg)
}
