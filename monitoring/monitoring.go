package monitoring

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "Rest",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

func Monitor() {
	port := ":8180"
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	fmt.Println("Monitoring port:", port)
	http.ListenAndServe(port, mux)
}

func RequestMonitoring(d time.Duration, status string) {
	requestMetrics.WithLabelValues(status).Observe(d.Seconds())
}
