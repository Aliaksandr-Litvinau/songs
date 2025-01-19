package middleware

import (
	"songs/internal/app/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware adds a middleware to collect Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.Request.URL.Path
		method := c.Request.Method

		// Write the metrics to Prometheus
		metrics.RequestDuration.WithLabelValues(path, method, status).Observe(duration)

		// Increment the total requests counter
		metrics.TotalRequests.WithLabelValues(path, method, status).Inc()
	}
}
