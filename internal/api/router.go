package api

import (
	"embed"
	"github.com/gin-gonic/gin"
	_ "github.com/kerim-dauren/user-service/internal/api/docs"
	"github.com/kerim-dauren/user-service/internal/api/http/middlewares"
	"github.com/kerim-dauren/user-service/internal/api/http/v1"
	"github.com/kerim-dauren/user-service/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

//go:embed docs/*
var swaggerDocs embed.FS

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "path", "result"},
	)
)

func init() {
	prometheus.MustRegister(requestDuration)
}

type RouterDeps struct {
	UserService domain.UserService
}

func NewHttpRouter(deps *RouterDeps) *gin.Engine {
	router := gin.New()

	if gin.Mode() == gin.DebugMode {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())

	// Swagger
	if gin.Mode() != gin.ReleaseMode {
		// Publish static Swagger files
		router.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Embedded Swagger files (accessible via /swagger/docs/*filepath)
		router.GET("/swagger/docs/*filepath", func(c *gin.Context) {
			http.FileServer(http.FS(swaggerDocs)).ServeHTTP(c.Writer, c.Request)
		})
	}

	// HealthCheck
	router.GET("/health", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Prometheus
	router.GET("/metrics", metricsHandler())

	// Routers
	apiV1 := router.Group("/api/v1")
	{
		apiV1.Use(
			middlewares.PrometheusMiddleware(requestDuration),
			//middlewares.TraceID(), //TODO tracing requests
		)

		userHandler := v1.NewUserHandler(deps.UserService)

		apiV1.POST("/users", userHandler.CreateUser)
		apiV1.GET("/users/:id", userHandler.GetUser)
		apiV1.PUT("/users/:id", userHandler.UpdateUser)
		apiV1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	return router
}

func metricsHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		Registry: prometheus.DefaultRegisterer,
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
