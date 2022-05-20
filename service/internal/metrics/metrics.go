package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var (
	requestsMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests",
			Help: "all requests by type",
		},
		[]string{"type", "resp_type"})

	requestsTimeSum = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_time_sum",
			Help: "sum of time of all requests",
		},
		[]string{"type", "resp_type"})

	cacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kitime_cache_hits",
			Help: "",
		},
		[]string{"is_hit"})
)

func Init() {
	prometheus.MustRegister(requestsMetric, requestsTimeSum, cacheHits)
}

func ObserveRequests(requestType, respType string) {
	requestsMetric.With(prometheus.Labels{
		"type":      requestType,
		"resp_type": respType,
	}).Inc()
}

func ObserveRequestsTimeSum(requestType, respType string, time float64) {
	requestsTimeSum.With(prometheus.Labels{
		"type":      requestType,
		"resp_type": respType,
	}).Add(time)
}

func ObserveCacheHits(isHit bool) {
	cacheHits.With(prometheus.Labels{
		"is_hit": strconv.FormatBool(isHit),
	}).Inc()
}
