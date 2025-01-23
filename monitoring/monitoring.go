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
	// Создаем метрику
	Namespace:  "Rest",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // Мапа с ключ значением и квантилями
}, []string{"status"}) // Просто статус заявки

func Monitor() {
	port := ":8180" // На отдельном портез апускаем стату монотринга
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler()) // Обязательно metrics
	fmt.Println("Monitoring port:", port)
	http.ListenAndServe(port, mux) // Запускаем
}

func RequestMonitoring(d time.Duration, status string) {
	requestMetrics.WithLabelValues(status).Observe(d.Seconds()) // Записываем метрику
}
