package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	HeaderTraceID     = "X-Trace-ID" // or use X-Request-ID for tracing requests ;)
	ContextTraceIDKey = "traceID"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(HeaderTraceID)
		if traceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("'%s' header is required", HeaderTraceID)})
			c.Abort()
			return
		}

		// Set the TraceID in the request context
		ctx := context.WithValue(c.Request.Context(), ContextTraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
