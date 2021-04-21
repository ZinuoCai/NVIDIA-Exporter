package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddr       = flag.String("web.listen-port", "9114", "An port to listen on for web interface and telemetry.")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
)

func main() {
	flag.Parse()

	metrics := NewMetrics()
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter for GPU Information Per Process</title></head>
			<body>
			<h1>A Prometheus Exporter for GPU Information Per Process</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting Server at http://localhost:%s%s", *listenAddr, *metricsPath)
	log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))
}
