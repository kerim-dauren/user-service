package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"time"
)

// PrometheusMiddleware is a middleware to collect Prometheus metrics for request durations.
// It uses a HistogramVec with labels: handler_name, method, and result (success/error).
func PrometheusMiddleware(requestDuration *prometheus.HistogramVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		elapsed := time.Since(start)

		result := "success"
		status := c.Writer.Status()
		if status >= 400 && status < 500 {
			result = "failed"
		}
		if status >= 500 {
			result = "error"
		}

		handler := parseHandlerName(c.HandlerName())
		path := c.FullPath()

		requestDuration.WithLabelValues(handler, path, result).Observe(elapsed.Seconds())
	}
}

// parseHandlerName extracts the short handler name from the full Gin handler signature.
func parseHandlerName(fullHandlerName string) string {
	// Example fullHandlerName:
	// "github.com/kerim-dauren/user-service/api/http/v1/routes.(*userHandler).createUser"

	// Extract the part inside the parentheses (e.g., "(*userHandler).createUser")
	start := strings.Index(fullHandlerName, "(*")
	if start == -1 {
		return fullHandlerName // Return as is if format is unexpected
	}
	end := strings.Index(fullHandlerName[start:], ")")
	if end == -1 {
		return fullHandlerName
	}

	// Trim to get "userHandler.createUser"
	trimmed := fullHandlerName[start+2 : start+end]
	parts := strings.Split(trimmed, ".")
	if len(parts) == 2 {
		// Return "govHandler" (handler name) and "getGbdflInfo" (path)
		return parts[0] + "." + parts[1]
	}
	return trimmed
}
