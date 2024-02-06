package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestsRPS = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_rps",
			Help:    "HTTP requests per second",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(HttpRequestsRPS)
}
