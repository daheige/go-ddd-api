package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/daheige/go-ddd-api/internal/infras/utils"
	"github.com/go-god/gutils"
)

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
				log.Println("error: ", err)
				log.Println("exec panic error", map[string]interface{}{
					"trace_error": string(debug.Stack()),
				})

				// services error
				http.Error(w, "services error!", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// AccessLog access log
func AccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println("exec begin", nil)
		log.Println("request before")
		log.Println("request method: ", r.Method)
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
			"exec_time": time.Since(start).Seconds(),
		})
	})
}
