package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	labelUrl           = "url"
	labelMethod        = "method"
	labelHost          = "host"
	labelStatusCode    = "statusCode"
	serverSystemName   = "httprelayserver"
	proxyletSystemName = "httprelaylet"
)

var (
	// Number of request hitting the system
	totalRequests prometheus.Counter

	// Number of failed  requests (histogram by error code)
	totalFailedRequests *prometheus.CounterVec
)

func initCommonMetrics(systemName string) {
	totalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_total_requests", systemName),
		Help: "The total number of request hitting the system",
	})

	totalFailedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_failed_requests", systemName),
		Help: "histogram of all the failed requests",
	}, []string{labelStatusCode})
}

func InitProxyServerMetrics() {
	initCommonMetrics(serverSystemName)
}

func InitProxyletMetrics() {
	initCommonMetrics(proxyletSystemName)

}

func IncTotalRequests() {
	if totalRequests == nil {
		return
	}

	totalRequests.Inc()
}

func IncTotalFailedRequests(statusCode string) {
	if totalFailedRequests == nil {
		return
	}

	totalFailedRequests.With(
		prometheus.Labels{
			labelStatusCode: statusCode,
		},
	).Inc()
}
