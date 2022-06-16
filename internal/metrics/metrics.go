package metrics

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/homework3/notification/internal/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func CreateMetricsServer(addr string, cfg *config.Config) *http.Server {
	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return metricsServer
}
