package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests",
			Help: "all requests by type",
		},
		[]string{"type", "resp_type"})

	RequestsTimeSum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_time_sum",
			Help: "sum of time of all requests",
		},
		[]string{"type", "resp_type"})
)

func Init() {
	prometheus.MustRegister(RequestsMetric, RequestsTimeSum)
}
