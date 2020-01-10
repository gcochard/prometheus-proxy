package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
)

// CollectorFactory should spawn custom Collector instances for target
type CollectorFactory func(target string) prometheus.Collector

type metricsHandler struct {
	factory CollectorFactory
}

func (m *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the intended target
	queries := r.URL.Query()
	target := queries.Get("target")
	if len(queries["target"]) != 1 || target == "" {
		http.Error(w, "'target' parameter must be specified once", 400)
		return
	}

	// Create a registry and the Collector
	reg := prometheus.NewRegistry()
	reg.MustRegister(m.factory(target))

	// Delegate to the default prometheus handler
	defaultHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	defaultHandler.ServeHTTP(w, r)
}

// ListenAndServe starts the HTTP server on address and port and listens for requests
// on /metrics?target=foo, passing foo to factory and starting a collection.
func ListenAndServe(address string, port string, factory CollectorFactory) error {
	serveAddr := net.JoinHostPort(address, port)
	log.Println("Exposing metrics on", serveAddr)

	handler := metricsHandler{factory: factory}
	http.Handle("/metrics", &handler)

	return http.ListenAndServe(serveAddr, nil)
}
