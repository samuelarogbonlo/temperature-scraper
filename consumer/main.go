package main

import (
	"context"
	"encoding/json"
	"net/http"

	"dia/interview/consumer/pkg/database"
	"dia/interview/consumer/pkg/kafkaConsumer"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Prometheus metrics custom config
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

const ConsumerPort = ":4002"

func main() {

	// Initialize database connection
	db, err := database.New()
	if err != nil {
		kafkaConsumer.ErrorLog.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize database schema
	if err := db.InitializeSchema(); err != nil {
		kafkaConsumer.ErrorLog.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Initialize Kafka consumer
	store := kafkaConsumer.NewDataStore()
	consumer := kafkaConsumer.NewConsumer(store, db)
	kafkaConsumer.InfoLog.Println("Kafka Consumer initialized successfully")

	ctx, cancel := context.WithCancel(context.Background())
	go kafkaConsumer.SetupConsumerGroup(ctx, consumer)
	defer cancel()
	kafkaConsumer.InfoLog.Println("Kafka Consumer Group initialized successfully")

	gin.SetMode(gin.ReleaseMode)
	// Setup Gin HTTP server
	router := gin.Default()

	// Middleware to collect metrics
	router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(path))
		c.Next() // Process request
		httpRequestsTotal.WithLabelValues(path).Inc()
		timer.ObserveDuration()
	})

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Endpoint to get data for a specific city
	router.GET("/city/:city", func(c *gin.Context) {
		city := c.Param("city")

		data, err := db.GetCityTemperature(city)
		if err != nil {
			kafkaConsumer.ErrorLog.Printf("Error querying database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if len(data) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the specified city"})
			return
		}

		prettyJSON, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			kafkaConsumer.ErrorLog.Printf("Error formatting JSON: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.Data(http.StatusOK, "application/json", prettyJSON)
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	kafkaConsumer.InfoLog.Printf("Starting HTTP server on port %s", ConsumerPort)
	kafkaConsumer.ErrorLog.Fatal(http.ListenAndServe(ConsumerPort, router))
}
