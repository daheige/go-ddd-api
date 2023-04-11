package middleware

import (
	"context"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/daheige/go-ddd-api/internal/infras/utils"
	"github.com/go-god/gutils"
	"github.com/go-god/logger"
	"go.uber.org/zap"
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
				// log.Println("error: ", err)
				// log.Println("exec panic error", map[string]interface{}{
				// 	"trace_error": string(debug.Stack()),
				// })

				ctx := r.Context()
				logger.Info(ctx, "exec panic error",
					zap.String("module", "web"), zap.String("trace_error", string(debug.Stack())),
				)

				if isBrokenPipe(ctx, err) {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				// services error
				http.Error(w, "services error", http.StatusInternalServerError)
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
		// log.Println("exec begin", nil)
		// log.Println("request before")
		// log.Println("request method: ", r.Method)
		// log.Println("request uri: ", r.RequestURI)

		// x-request-id
		reqId := r.Header.Get("x-request-id")
		if reqId == "" {
			reqId = gutils.Uuid()
		}

		// log.Println("log_id: ", reqId)
		userAgentKey := logger.CtxKey{Name: "user-agent"}
		userAgent := r.Header.Get("User-Agent")
		r = utils.ContextSet(r, logger.XRequestID, reqId)
		r = utils.ContextSet(r, logger.ReqClientIP, r.RemoteAddr)
		r = utils.ContextSet(r, logger.RequestMethod, r.Method)
		r = utils.ContextSet(r, logger.RequestURI, r.RequestURI)
		r = utils.ContextSet(r, userAgentKey, userAgent)

		logger.Info(r.Context(), "exec begin",
			zap.String("module", "web"), zap.String(userAgentKey.String(), userAgent),
		)

		h.ServeHTTP(w, r)

		// log.Println("exec end", map[string]interface{}{
		// 	"exec_time": time.Since(start).Seconds(),
		// })

		logger.Info(r.Context(), "exec end", map[string]interface{}{
			"exec_time": time.Since(start).Seconds(),
		})
	})
}

func isBrokenPipe(ctx context.Context, err interface{}) bool {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, exist := ne.Err.(*os.SyscallError); exist {
			errMsg := strings.ToLower(se.Error())
			// logger error
			logger.Error(ctx, "os syscall error", map[string]interface{}{
				"trace_error": errMsg,
			})

			if strings.Contains(errMsg, "broken pipe") ||
				strings.Contains(errMsg, "reset by peer") ||
				strings.Contains(errMsg, "request headers: small read buffer") ||
				strings.Contains(errMsg, "unexpected EOF") ||
				strings.Contains(errMsg, "i/o timeout") {
				brokenPipe = true
			}
		}
	}

	return brokenPipe
}
