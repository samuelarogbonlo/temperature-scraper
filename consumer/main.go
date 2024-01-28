package main

import (
	"context"
	"log"
	"net/http"
	_ "os" // Import os for environment variables

	"dia/interview/consumer/pkg/database"      // Import your database package
	"dia/interview/consumer/pkg/kafkaConsumer" // Adjust the import path to your kafkaConsumer package

	"github.com/gin-gonic/gin"
)

const ConsumerPort = ":8081"

func main() {
    // Fetch database connection string from environment variable
    dbConnString := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

    // if dbConnString == "" {
    //     log.Fatal("DB_CONN_STRING environment variable is required")
    // }

    // Initialize database connection
    db, err := database.New(dbConnString)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize database schema
    if err := db.InitializeSchema(); err != nil {
        log.Fatalf("Failed to initialize database schema: %v", err)
    }

    // Initialize Kafka consumer
    store := kafkaConsumer.NewDataStore()
    consumer := kafkaConsumer.NewConsumer(store, db)

    ctx, cancel := context.WithCancel(context.Background())
    go kafkaConsumer.SetupConsumerGroup(ctx, consumer) // Pass reference of consumer
    defer cancel()

    // Setup Gin HTTP server
    router := gin.Default()
    router.GET("/city-temperature", func(c *gin.Context) {
        data := store.GetAll()
        log.Printf("HTTP request: returning %d records", len(data)) // Add logging
        c.JSON(http.StatusOK, gin.H{"data": data})
    })

    log.Printf("Starting HTTP server on port %s", ConsumerPort)
    log.Fatal(http.ListenAndServe(ConsumerPort, router))
}
