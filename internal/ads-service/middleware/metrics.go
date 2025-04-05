package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func IncreaseTotalRequestsMetric(metric *prometheus.CounterVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		metric.WithLabelValues("total_requests").Inc()
	}
}
