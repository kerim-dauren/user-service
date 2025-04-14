package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTraceID(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("ValidTraceID", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(TraceID())

		// Add a test handler to verify context propagation
		var receivedTraceID string
		router.GET("/test", func(c *gin.Context) {
			// Extract the traceID from the request context
			traceID, exists := c.Request.Context().Value(ContextTraceIDKey).(string)
			if exists {
				receivedTraceID = traceID
			}
			c.Status(http.StatusOK)
		})

		// Create request with valid trace ID
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		expectedTraceID := "test-trace-123"
		req.Header.Set(HeaderTraceID, expectedTraceID)

		// Perform the request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedTraceID, receivedTraceID, "TraceID should be passed to handler context")
	})

	t.Run("MissingTraceID", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(TraceID())

		// Add handler that should not be called
		router.GET("/test", func(c *gin.Context) {
			// This should not be called when trace ID is missing
			t.Error("Handler was called despite missing TraceID")
			c.Status(http.StatusOK)
		})

		// Create request without trace ID
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Perform the request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), HeaderTraceID)
	})

	t.Run("EmptyTraceID", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(TraceID())

		// Add handler that should not be called
		router.GET("/test", func(c *gin.Context) {
			// This should not be called when trace ID is empty
			t.Error("Handler was called despite empty TraceID")
			c.Status(http.StatusOK)
		})

		// Create request with empty trace ID
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(HeaderTraceID, "")

		// Perform the request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), HeaderTraceID)
	})

	t.Run("ContextPropagation", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(TraceID())

		// Add a chain of handlers to verify context propagation through multiple handlers
		propagationVerified := false

		router.GET("/test", func(c *gin.Context) {
			// First middleware
			traceID, exists := c.Request.Context().Value(ContextTraceIDKey).(string)
			assert.True(t, exists)
			assert.Equal(t, "propagation-test", traceID)
			c.Next()
		}, func(c *gin.Context) {
			// Second middleware - should still have access to the same context
			traceID, exists := c.Request.Context().Value(ContextTraceIDKey).(string)
			assert.True(t, exists)
			assert.Equal(t, "propagation-test", traceID)
			propagationVerified = true
			c.Status(http.StatusOK)
		})

		// Create request with valid trace ID
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(HeaderTraceID, "propagation-test")

		// Perform the request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, propagationVerified, "Context propagation was not verified")
	})

	t.Run("MultipleRequests", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(TraceID())

		var receivedTraceIDs []string
		router.GET("/test", func(c *gin.Context) {
			traceID, _ := c.Request.Context().Value(ContextTraceIDKey).(string)
			receivedTraceIDs = append(receivedTraceIDs, traceID)
			c.Status(http.StatusOK)
		})

		// Multiple requests with different trace IDs
		traceIDs := []string{"trace-1", "trace-2", "trace-3"}

		for _, id := range traceIDs {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set(HeaderTraceID, id)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Assert all trace IDs were received correctly
		assert.Equal(t, traceIDs, receivedTraceIDs)
	})
}

func TestTraceIDExtraction(t *testing.T) {
	// Test the utility function for extracting trace ID from context
	t.Run("ExtractTraceIDFromContext", func(t *testing.T) {
		// Create a context with a trace ID
		expectedTraceID := "extract-test-id"
		ctx := context.WithValue(context.Background(), ContextTraceIDKey, expectedTraceID)

		// Extract using a helper function (this would typically be defined in your application)
		extractTraceID := func(ctx context.Context) (string, bool) {
			id, ok := ctx.Value(ContextTraceIDKey).(string)
			return id, ok
		}

		traceID, exists := extractTraceID(ctx)

		assert.True(t, exists)
		assert.Equal(t, expectedTraceID, traceID)
	})
}
