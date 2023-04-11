package monitor

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/daheige/go-ddd-api/internal/infras/config"
	"github.com/go-god/monitor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusMonitor service prometheus monitor
type PrometheusMonitor struct {
	AppConfig *config.AppConfig `inject:""`
}

// Monitor prometheus and pprof monitor
func (s *PrometheusMonitor) Monitor() {
	// Registration monitoring index
	// Performance monitoring for web applications
	// if it is a job/rpc service, does not require these two lines
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)

	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	go func() {
		PProfAddress := fmt.Sprintf("0.0.0.0:%d", s.AppConfig.PProfPort)
		log.Printf("server pprof run on:%s\n", PProfAddress)

		httpMux := http.NewServeMux() // create a http ServeMux instance
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", s.check)

		// metrics monitor
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(PProfAddress, httpMux); err != nil {
			log.Println(err)
		}
	}()
}

// check PProf Heartbeat detection
func (s *PrometheusMonitor) check(w http.ResponseWriter, r *http.Request) {
	log.Println("request_uri: ", r.RequestURI)
	w.Write([]byte(`{"code": 0,"message":"ok"}`))
}
