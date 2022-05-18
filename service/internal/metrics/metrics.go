package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests",
			Help: "all requests by type",
		},
		[]string{"type", "resp_type", "time"})
)

func Init() {
	prometheus.MustRegister(RequestsMetric)
}
