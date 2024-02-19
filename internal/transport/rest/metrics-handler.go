package rest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type respWritterImpl struct {
	http.ResponseWriter
	statusCode int
}

func newRespWriterImpl(w http.ResponseWriter) *respWritterImpl {
	return &respWritterImpl{w, http.StatusOK}
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := newRespWriterImpl(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

// Function for handle getting metrics logic
func HandleMetricsLogic(mux *mux.Router) {
	// Register custom metrics
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	// Use custom MiddlewareFunc
	mux.Use(prometheusMiddleware)

	// Prometheus endpoint
	mux.Path("/metrics").Handler(promhttp.Handler())
}
