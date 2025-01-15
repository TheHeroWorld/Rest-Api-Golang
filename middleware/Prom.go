package middleware

import (
	"My_Frist_Golang/monitoring"
	"net/http"
	"time"
)

func MonitorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Мидлвейр мониторинга времени обработки запросов и статуса
		start := time.Now() // Иницилизурем время начала запроса

		defer func() {
			monitoring.RequestMonitoring(time.Since(start), r.Method) // Вызываем функцию в отложенном режиме куда оптравялем разинцу старта и сейчас, а так-же метод
		}()
		next.ServeHTTP(w, r)
	})
}
