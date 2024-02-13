package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const Namespace = "golang_api"

type PrometheusMetrics struct {
	UserMetrics     prometheus.Gauge
	RequestsMetrics *prometheus.CounterVec
}

var metrics *PrometheusMetrics = nil

func InitPrometheusMetrics() *PrometheusMetrics {
	if metrics == nil {
		metrics = &PrometheusMetrics{
			UserMetrics: promauto.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "users_total",
				Help:      "Total users",
			}),
			RequestsMetrics: promauto.NewCounterVec(prometheus.CounterOpts{
				Namespace: Namespace,
				Name:      "requests_total",
				Help:      "Total http requests",
			}, []string{"code"}),
		}
	}

	return metrics
}
